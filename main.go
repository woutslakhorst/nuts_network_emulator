package main

func main() {
	network := NewNetwork(
		Config{
			Nodes:    5,
			Rate:     30.0,
			Variance: 0.2,
			Rounds:   1000,
		},
	)

	// block call
	if err := network.Start(); err != nil {
		print(err)
	}
}
