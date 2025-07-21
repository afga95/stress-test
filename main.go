package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Estrutura para armazenar resultados de cada request
type RequestResult struct {
	StatusCode int
	Error      error
	Duration   time.Duration
}

// Estrutura para o relatório final
type Report struct {
	TotalTime       time.Duration
	TotalRequests   int
	SuccessRequests int
	StatusCodes     map[int]int
	ErrorCount      int
	AverageResponse time.Duration
}

func main() {
	// Definição dos parâmetros da linha de comando
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	// Validação dos parâmetros
	if *url == "" {
		fmt.Println("Erro: URL é obrigatória")
		fmt.Println("Uso: programa --url=<URL> --requests=<número> --concurrency=<número>")
		os.Exit(1)
	}

	if *requests <= 0 {
		fmt.Println("Erro: Número de requests deve ser maior que 0")
		os.Exit(1)
	}

	if *concurrency <= 0 {
		fmt.Println("Erro: Nível de concorrência deve ser maior que 0")
		os.Exit(1)
	}

	if *concurrency > *requests {
		*concurrency = *requests
	}

	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total de requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Println("==================================================")

	// Executar o teste de carga
	report := runLoadTest(*url, *requests, *concurrency)

	// Exibir relatório
	printReport(report)
}

func runLoadTest(url string, totalRequests, concurrency int) Report {
	startTime := time.Now()

	// Canal para receber resultados
	results := make(chan RequestResult, totalRequests)

	// Canal para controlar o número de requests
	requestChan := make(chan int, totalRequests)

	// Preencher o canal com os números dos requests
	for i := 0; i < totalRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// WaitGroup para aguardar todas as goroutines
	var wg sync.WaitGroup

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Iniciar workers concorrentes
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(client, url, requestChan, results, &wg)
	}

	// Aguardar todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Processar resultados
	report := Report{
		StatusCodes: make(map[int]int),
	}

	var totalDuration time.Duration

	for result := range results {
		report.TotalRequests++
		totalDuration += result.Duration

		if result.Error != nil {
			report.ErrorCount++
		} else {
			report.StatusCodes[result.StatusCode]++
			if result.StatusCode == 200 {
				report.SuccessRequests++
			}
		}
	}

	report.TotalTime = time.Since(startTime)
	if report.TotalRequests > 0 {
		report.AverageResponse = totalDuration / time.Duration(report.TotalRequests)
	}

	return report
}

func worker(client *http.Client, url string, requestChan <-chan int, results chan<- RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for range requestChan {
		result := makeRequest(client, url)
		results <- result
	}
}

func makeRequest(client *http.Client, url string) RequestResult {
	start := time.Now()

	resp, err := client.Get(url)
	duration := time.Since(start)

	result := RequestResult{
		Duration: duration,
		Error:    err,
	}

	if err != nil {
		return result
	}

	result.StatusCode = resp.StatusCode
	resp.Body.Close()

	return result
}

func printReport(report Report) {
	fmt.Println("\n" + "==================================================")
	fmt.Println("RELATÓRIO DE TESTE DE CARGA")
	fmt.Println("==================================================")

	fmt.Printf("Tempo total de execução: %v\n", report.TotalTime)
	fmt.Printf("Total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", report.SuccessRequests)
	fmt.Printf("Requests com erro: %d\n", report.ErrorCount)
	fmt.Printf("Tempo médio de resposta: %v\n", report.AverageResponse)

	if len(report.StatusCodes) > 0 {
		fmt.Println("\nDistribuição de códigos de status:")
		for statusCode, count := range report.StatusCodes {
			percentage := float64(count) / float64(report.TotalRequests) * 100
			fmt.Printf("  %d: %d (%.1f%%)\n", statusCode, count, percentage)
		}
	}

	// Calcular taxa de sucesso
	successRate := float64(report.SuccessRequests) / float64(report.TotalRequests) * 100
	fmt.Printf("\nTaxa de sucesso: %.1f%%\n", successRate)

	// Calcular requests por segundo
	requestsPerSecond := float64(report.TotalRequests) / report.TotalTime.Seconds()
	fmt.Printf("Requests por segundo: %.2f\n", requestsPerSecond)

	fmt.Println("==================================================")
}
