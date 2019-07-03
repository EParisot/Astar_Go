package main

func (env *Env) botPlayer(algo string) {
	// Select Algo
	//TODO append algos here
	switch {
	case algo == "Astar":
		go env.aStar()
	}
}

func (env *Env) aStar() {

}
