package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

// Imprime los detalles puntero, goleador y asistente del link pasado
func printDetails(link string, wg *sync.WaitGroup) {
	//Avisar al waiter que se termino la ejecución
	defer wg.Done()

	//Variables que se utilizarán
	var (
		titulo       string = "No encontrado"
		puntero      string = "No encontrado"
		goleador     string = "No encontrado"
		asistente    string = "No encontrado"
		encontroAlgo bool   = false
	)

	details := colly.NewCollector()

	//Encuentra el título de la sección
	details.OnHTML("#titulos", func(h *colly.HTMLElement) {
		titulo = h.Text
	})

	//Encuentra el puntero del link (equipo que va primero)
	details.OnHTML("#posiciones > tbody > tr.punt > td[align]", func(h *colly.HTMLElement) {
		puntero = h.Text
		encontroAlgo = true
	})

	// Encuentra el máximo asistente/goleado
	details.OnHTML("#goleadorest", func(h *colly.HTMLElement) {
		encontroAlgo = true
		if h.ChildText("tbody > tr:nth-child(1) > th:nth-child(2)") == "Asist." {
			asistente = h.ChildText("tbody > tr.punt > td:nth-child(1)")
		} else {
			goleador = h.ChildText("tbody > tr.punt > td:nth-child(1)")

		}
	})

	// Visita el link
	details.Visit(link)

	// Imprime los datos si es que encontró al menos uno de ellos
	if encontroAlgo {
		fmt.Println(titulo)
		fmt.Println("Puntero: ", puntero)
		fmt.Println("Máximo goleador: ", goleador)
		fmt.Println("Máximo asistente: ", asistente)
		fmt.Println()
	}

}

func main() {
	// Waiter que esperará las funcinoes concurrentes
	var wg sync.WaitGroup

	c := colly.NewCollector()

	// Visita cada link que encuentra en la barra del menu
	c.OnHTML("#accordian > ul > li > ul.items-menu > li > a", func(h *colly.HTMLElement) {
		wg.Add(1)

		// Link a visitar
		link := h.Attr("href")
		link = h.Request.AbsoluteURL(link)

		// Imprime los daots del link
		go printDetails(link, &wg)

	})

	c.Visit("https://www.promiedos.com.ar/")

	// Espera que acaben todas las funciones que se ejcutan de manera concurrente
	wg.Wait()

}
