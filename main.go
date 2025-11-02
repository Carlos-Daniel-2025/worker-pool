package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//Declaramos la función worker el cual procesara las tareas haciendo uso del canal jobs
func worker(id int, ctx context.Context, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	//Aviso que para cuando el worker acabe
	defer wg.Done() 

	for {
		select {
		//Colocamos un case para recibir una tarea del canal jobs
		case job, ok := <-jobs:
			if !ok {
				//Colocamos un return para en dado caso el canal se encuntre cerrado y ya no haya más trabajos
				return
			}

			fmt.Println("Worker", id, "started job", job)

			//Hacemos uso del time.Sleep para simular trabajos con 1s de espera
			time.Sleep(time.Second)

			//Mensaje de tarea finalizada
			fmt.Println("Worker", id, "finished job", job)

			//Se envian los resultados al canal results
			results <- job * 2

		//Colocamos un case por si el contexto recibe una señal de cancelación este detenga el worker
		case <-ctx.Done():
			//Mensaje que se mostrara al usuario
			fmt.Println("Worker", id, "detenido por cancelación del contexto")
			return
		}
	}
}

func main() {
	//Hacemos uso del time.now para marcar el inicio del tiempo de ejecución
	start := time.Now()

	//Definimos cuantos workers tendre, en este caso 4, y cuántos jobs hay, en este caso seran 20
	numWorkers := 4
	numJobs := 20

	//Colocamos un contexto con cancelación manual
	ctx, cancel := context.WithCancel(context.Background())

	//Creamos los canales de jobs y de resultados
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	//Declaramos la variable deWaitGroup para esperar a los workers, igualmente la colocamos en la 
	//parte del inicio de imports para que sea reconocida haciendo uso de "sync"
	var wg sync.WaitGroup

	//Mando llamar los workers
	for i := 1; i <= numWorkers; i++ {
		//Agregamos un worker al contador del WaitGroup
		wg.Add(1)
		go worker(i, ctx, jobs, results, &wg)
	}

	//Envíamos los trabajos al canal jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	//Hacemos uso de un close para decirle al codigo que finalice el canal de jobs para indicar que no hay más
	close(jobs) 

	//Esperamos a que todos los workers se terminen de ejecutar
	wg.Wait()

	//Cancelamos el contexto ctx, ya que no lo ocuparemos
	cancel()

	//Imprimimos todos los resultados mandados
	for r := 1; r <= numJobs; r++ {
		fmt.Println("Result:", <-results)
	}

	//Calculmos tiempo total e imprimimos los resultados en pantalla
	//Usamos elapsed para mostrar el tiempo transcurrido
	elapsed := time.Since(start)
	//Mensaje que nos dice que todas las trabajos fueron realizados
	fmt.Println("\nAll jobs processed")
	//Mensaje que nos proporciona el tiempo total
	fmt.Println("\nTiempo total de ejecución:", elapsed)
	fmt.Println()
}