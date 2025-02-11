package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Route struct {
	StartingStation string
	EndingStation   string
	NumberOfStops   int
	RouteLength     float64
}

func (r Route) String() string {
	return fmt.Sprintf("Route: %s -> %s, Stops: %d, Route Length: %.2f km",
		r.StartingStation, r.EndingStation, r.NumberOfStops, r.RouteLength)
}

func readRoutesFromFile(filename string) ([]Route, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var routes []Route
	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		stopCount, err := strconv.Atoi(record[2])
		if err != nil || stopCount < 0 {
			continue
		}
		lengthKm, err := strconv.ParseFloat(record[3], 64)
		if err != nil || lengthKm <= 0 {
			continue
		}

		routes = append(routes, Route{
			StartingStation: record[0],
			EndingStation:   record[1],
			NumberOfStops:   stopCount,
			RouteLength:     lengthKm,
		})
	}

	return routes, nil
}

func sortRoutesByLength(routes []Route) {
	for i := 0; i < len(routes)-1; i++ {
		for j := 0; j < len(routes)-i-1; j++ {
			if routes[j].RouteLength > routes[j+1].RouteLength {
				routes[j], routes[j+1] = routes[j+1], routes[j]
			}
		}
	}
}

func countRoutesWithAverageDistance(routes []Route, x float64) int {
	count := 0
	for _, route := range routes {
		if route.NumberOfStops >= 0 {
			averageDistance := route.RouteLength / float64(route.NumberOfStops+1)
			if averageDistance < x {
				count++
			}
		}
	}
	return count
}

func getRoutesStartingFrom(routes []Route, station string) []Route {
	var filteredRoutes []Route
	for _, route := range routes {
		if strings.EqualFold(route.StartingStation, station) {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	return filteredRoutes
}

func getRoutesWithMaxStops(routes []Route) []Route {
	if len(routes) == 0 {
		return nil
	}

	maxStops := routes[0].NumberOfStops
	for _, route := range routes {
		if route.NumberOfStops > maxStops {
			maxStops = route.NumberOfStops
		}
	}

	var maxStopRoutes []Route
	for _, route := range routes {
		if route.NumberOfStops == maxStops {
			maxStopRoutes = append(maxStopRoutes, route)
		}
	}

	return maxStopRoutes
}

func main() {
	routes, err := readRoutesFromFile("routes.csv")
	if err != nil {
		fmt.Println("Error reading routes:", err)
		return
	}

	// Вивід всіх маршрутів з файлу
	fmt.Println("\nAll routes:")
	for _, route := range routes {
		fmt.Printf("%s\n", route)
	}

	// Сортування за протяжністю
	sortRoutesByLength(routes)
	fmt.Println("\nRoutes sorted by distance:")
	for _, route := range routes {
		fmt.Printf("%s\n", route)
	}

	// Кількість маршрутів, для яких середня довжина між зупинками менша за Х
	X := 15.0
	count := countRoutesWithAverageDistance(routes, X)
	fmt.Printf("\nNumber of routes with average length between stops less than %.2f kilometers: %d\n", X, count)

	// Список маршрутів, які починаються в станції Х
	station := "Dnipro"
	filteredRoutes := getRoutesStartingFrom(routes, station)
	fmt.Printf("\nRoutes that starts at the station %s:\n", station)
	for _, route := range filteredRoutes {
		fmt.Printf("%s\n", route)
	}

	// Маршрути з максимальною кількістю зупинок
	maxStopRoutes := getRoutesWithMaxStops(routes)
	fmt.Println("\nRoutes with the maximum number of stops:")
	for _, route := range maxStopRoutes {
		fmt.Printf("%s\n", route)
	}
}
