/*Huey Padua
Programming Assignment 1
Scheduling algorithms in GO
Due Date: July 15, 2018

I, Huey (HU658731), affirm that this program is entirely
my own work and that I have neither developed my code
together with any other person, nor copied any code from
any other person, nor permitted my code to be copied or otherwise
used by any other person, nor have I copied, modified, or otherwise
used programs created by others. I acknowledge that any violation of
the above terms will be treated as academic dishonesty.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Schedule struct {
	process_count int
	duration      int
	algorithm     string
	processes     []Process
	quantum       int
}

type Process struct {
	pid      string
	arrival  int
	burst    int
	selected int
}

func main() {
	//command line arguments for reading input file
	// and write to output file
	input := os.Args[1]

	var schedule Schedule
	var process Process

	infile, err := os.Open(input)
	defer infile.Close()

	if err != nil {
		fmt.Println(err)
	}

	// read file per line using scanner
	scanner := bufio.NewScanner(infile)
	//loop through each line
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, " ")
		if s[0] == "processcount" {

			ns := strings.Split(s[1], "\t")
			pc, err := strconv.Atoi(ns[0])

			if err != nil {
				fmt.Println(err)
			}
			schedule.process_count = pc
			//fmt.Println(schedule.proccess_count)
		}
		if s[0] == "runfor" {
			ns := strings.Split(s[1], "\t")
			dur, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.duration = dur
			//fmt.Println(schedule.duration)
		}
		if s[0] == "use" {
			ns := strings.Split(s[1], "\t")
			schedule.algorithm = ns[0]
			//fmt.Println(schedule.algorithm)
		}
		if s[0] == "quantum" {
			ns := strings.Split(s[1], "\t")
			quantum, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.quantum = quantum
			// fmt.Printf("%+v\n", schedule.quantum)
		}
		if s[0] == "process" {
			process.pid = s[2]
			a, err := strconv.Atoi(s[4])
			if err != nil {
				fmt.Println(err)
			}
			process.arrival = a
			b, err := strconv.Atoi(s[6])
			if err != nil {
				fmt.Println(err)
			}
			process.burst = b

			schedule.processes = append(schedule.processes, process)
			//fmt.Println(schedule.processes)
		}
		//fmt.Println(line)
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

	// printSchedule(&schedule)

	if schedule.algorithm == "fcfs" {
		fcfs(&schedule)
	} else if schedule.algorithm == "sjf" {
		sjf(&schedule)
	} else if schedule.algorithm == "rr" {
		rr(&schedule)
	}
}

//Helper function to print schedule information
func printSchedule(schedule *Schedule) {
	fmt.Println()
	fmt.Println("Algorithm: ", schedule.algorithm)
	fmt.Printf("Quantum %+v\n", schedule.quantum)
	fmt.Println("Process count: ", schedule.process_count)
	fmt.Println("Duration: ", schedule.duration)

	for i := 0; i < schedule.process_count; i++ {
		fmt.Println("Processes: ", schedule.processes[i])
		//fmt.Println(schedule.processes[i].pid)
	}
}

//Function to run first come first serve algorithm
func fcfs(schedule *Schedule) {
	outfile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal("Cannot create file..", err)
	}
	defer outfile.Close()

	fmt.Fprintf(outfile, "%3d processes\n", schedule.process_count)
	fmt.Fprintf(outfile, "Using First-Come First-Served\n\n")

	count := schedule.process_count
	var processes []Process
	for a := 0; a < schedule.process_count; a++ {
		processes = append(processes, schedule.processes[a])
	}

	b := 0
	lock := 0
	idle := 1
	time := 0
	// Loop through entire schedule duration
	for i := 0; i < schedule.duration+1; i++ {

		//If all processors are finished, it will go into idle
		if count == 0 {
			fmt.Fprintf(outfile, "Time %3d : Idle\n", b)
			b++
			//end loop once schedule duration has been reached
			if b == schedule.duration {
				break
			}
		}
		// Loop through all processors
		for j := 0; j < schedule.process_count; j++ {

			if processes[j].arrival == i {
				fmt.Fprintf(outfile, "Time %3d : %+v arrived\n", i, processes[j].pid)
			}
			// Select processor based on arrival
			// and only if no processes are running
			if i == processes[j].arrival && lock == 0 && idle == 1 {
				fmt.Fprintf(outfile, "Time %3d : %+v selected (burst %3d)\n", i, processes[j].pid,
					processes[j].burst)

				b = time + processes[j].burst
				processes[j].selected = time
				idle = 0
				lock = 1
			} else if time == processes[j].arrival && lock == 1 {
				// continue if next process arrives but there's a processsor running
				i++
			}

			// Print processor that's finished if time is reached based on burst
			if time <= b && lock == 1 && idle == 0 {
				fmt.Fprintf(outfile, "Time %3d : %+v finished\n", b, processes[j].pid)
				lock = 0
				idle = 1
				b = time + processes[j].burst
				count--
				// Selects a new process to run if there's still any waiting
				// to be selected and finished
				if time == processes[j].arrival && lock == 0 && count != 0 {
					// fmt.Printf("Time %3d : %+v selected (burst %3d)\n", i, processes[j].pid,
					// 	processes[j].burst)
					// processes[j].selected = time
					// b = time + processes[j].burst
				} else if time == processes[j].arrival && lock == 1 {
					// Continue if a processor has arrived
					// but theres a processor already running
					i++
				}
			}
			// Goes into idle until another process arrives based on given time
			if time < processes[j].arrival && idle == 1 {
				time = i + b
				fmt.Fprintf(outfile, "Time %3d : Idle\n", time)

			}
		}
	}
	fmt.Fprintf(outfile, "Finished at time %3d\n\n", schedule.duration)

	// Prints out waiting and turnaround times for each process
	for i := 0; i < schedule.process_count; i++ {
		waitTime := processes[i].selected - processes[i].arrival
		turnaround := waitTime + processes[i].burst
		fmt.Fprintf(outfile, "%+v wait %3d turnaround %3d\n", processes[i].pid, waitTime, turnaround)
	}
}

// Function to run Shortest Job First algorithm
func sjf(schedule *Schedule) {
	outfile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal("Cannot create file..", err)
	}
	defer outfile.Close()

	fmt.Fprintf(outfile, "%3d processes\n", schedule.process_count)
	fmt.Fprintf(outfile, "Using preemptive Shortest Job First\n\n")

	var processes []Process
	for a := 0; a < schedule.process_count; a++ {
		processes = append(processes, schedule.processes[a])
	}

	count := schedule.process_count
	remainingBurst := make([]int, count)
	wait := make([]int, count)
	turnaround := make([]int, count)
	selected := -1

	// Original burst values are kept in order to
	// calculate wait and turnaround times
	for i := 0; i < count; i++ {
		remainingBurst[i] = processes[i].burst
	}

	for time := 0; time < schedule.duration; time++ {

		for i := 0; i < count; i++ {
			// If a process has arrived
			if processes[i].arrival == time {
				fmt.Fprintf(outfile, "Time %3d : %+v arrived\n", time, processes[i].pid)

				// If a process has not been selected
				// select new process
				if selected == -1 {
					selected = i
					fmt.Fprintf(outfile, "Time %3d : %+v selected (burst %3d)\n",
						time, processes[selected].pid, remainingBurst[selected])
				} else {
					// If a process has been selected,
					// check if the new available process has a shorter burst
					for j := 0; j < count; j++ {

						if remainingBurst[j] < remainingBurst[selected] && remainingBurst[j] > 0 {
							// If process has a shorter burst, select
							// based on preemptive sjf
							selected = j
							fmt.Fprintf(outfile, "Time %3d : %+v selected (burst %3d)\n",
								time, processes[selected].pid, remainingBurst[selected])
						}
					}
				}
			}
		}
		// If a process had been selected already
		if selected != -1 {
			// remainingBurst decremented
			remainingBurst[selected]--
			// If remainingBurst reaches 0,
			// process is finished
			if remainingBurst[selected] == 0 {
				// time is off by 1,
				// fixed issue to match output
				time = time + 1
				fmt.Fprintf(outfile, "Time %3d : %+v finished\n", time,
					processes[selected].pid)

				// wait and turnaround times calculated
				wait[selected] = time - (processes[selected].burst + processes[selected].arrival)
				turnaround[selected] = wait[selected] + processes[selected].burst

				// Variables reinitialized to find next process
				// with shortest burst
				shortestProcess := -1
				shortestBurst := math.MaxInt8

				// Selects the next shortest process to run
				// from the availabe processes
				for i := 0; i < count; i++ {

					if processes[i].arrival <= time && remainingBurst[i] > 0 {
						if remainingBurst[i] < shortestBurst {
							shortestProcess = i
							shortestBurst = remainingBurst[i]
						}
					}
				}
				selected = shortestProcess
				// If a new process has been selected, print desired outputs
				if selected != -1 {
					fmt.Fprintf(outfile, "Time %3d : %+v selected (burst %+v)\n", time,
						processes[selected].pid, remainingBurst[selected])
				} else {
					// Processor goes on Idle,
					// if a new a process has not been selected
					if time < schedule.duration {
						fmt.Fprintf(outfile, "Time %3d : Idle\n", time)
					}
				}
			}

		} else {
			// Processor goes on Idle when all the processes
			// has been finished
			fmt.Fprintf(outfile, "Time %3d : Idle\n", time)
		}
	}
	fmt.Fprintf(outfile, "Finished at time %3d\n\n", schedule.duration)

	for i := 0; i < count; i++ {
		fmt.Fprintf(outfile, "%+v wait %3d turnaround %3d\n", processes[i].pid,
			wait[i], turnaround[i])
	}
}

// Function to run Round-Robin algorithm
func rr(schedule *Schedule) {
	outfile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal("Cannot create file..", err)
	}
	defer outfile.Close()

	fmt.Fprintf(outfile, "%3d processes\n", schedule.process_count)
	fmt.Fprintf(outfile, "Using Round-Robin\n")
	fmt.Fprintf(outfile, "Quantum %3d\n\n", schedule.quantum)

	var processes []Process
	for a := 0; a < schedule.process_count; a++ {
		processes = append(processes, schedule.processes[a])
	}

	count := schedule.process_count
	quantum := schedule.quantum
	remainingBurst := make([]int, count)
	wait := make([]int, count)
	turnaround := make([]int, count)

	time := 0
	selected := -1
	last := -1
	finished := 0
	done := 0
	ready := 0

	// original burst values are kept in order
	// to calculate wait and turnaround
	for i := 0; i < count; i++ {
		remainingBurst[i] = processes[i].burst
	}

	// Loop until processes are all done
	for {
		for i := 0; i < count; i++ {
			if processes[i].arrival == time || processes[i].arrival > (time-quantum) && processes[i].arrival < time {
				fmt.Fprintf(outfile, "Time %3d : %+v arrived\n", processes[i].arrival, processes[i].pid)

				// If a process has not been selected, select the first
				if selected == -1 {
					selected = i
					ready = 1
				}
			}
		}
		if ready == 1 {
			// if there is still processes to run
			if finished != count {
				// If a process has not been selected
				// select the first available
				if selected != -1 {

					for i := 0; i < count; i++ {
						// If the process that is ready to switch has arrived
						if (processes[selected].arrival % count) <= time {

							if (remainingBurst[selected] % count) > 0 {

								selected = (selected) % count
							}
						}
					}
				} else {
					// Process is not currently selected,
					// loop thru and select the next process based on the last selected
					for i := 1; i <= count; i++ {

						if (processes[(last+i)%count].arrival) <= time {

							if (remainingBurst[(last+i)%count]) > 0 {

								selected = (last + i) % count
							}
						}
					}
				}
				// If a process is selected
				if selected != -1 {

					fmt.Fprintf(outfile, "Time %3d : %+v selected (burst %3d)\n", time, processes[selected].pid,
						remainingBurst[selected])
					if remainingBurst[selected]-quantum < 0 {

						time += remainingBurst[selected]
					} else {
						time += quantum
					}
					remainingBurst[selected] -= quantum
					last = selected

					if remainingBurst[selected] <= 0 {

						fmt.Fprintf(outfile, "Time %3d : %+v finished\n", time, processes[selected].pid)

						finished++
						// calculating wait and turnaround times
						wait[selected] = time - (processes[selected].burst + processes[selected].arrival)
						turnaround[selected] = wait[selected] + processes[selected].burst

						if finished == count && time < schedule.duration {
							fmt.Fprintf(outfile, "Time %3d : Idle\n", time)
							time++
						}
						selected = -1
					}
				} else {
					// Processor idles if a process has not been selected
					fmt.Fprintf(outfile, "Time %3d : Idle\n", time)
					time++
				}
			} else {
				// All processes finished,
				// processor idles till duration is complete
				fmt.Fprintf(outfile, "Time %3d : Idle\n", time)
				time++
			}
			// Once time reaches scheduled duration,
			// processes are complete
			if time >= schedule.duration {
				done = 1
				break
			}
		}
		// Stop loop once duration is reached
		if done == 1 {
			break
		}
	}
	fmt.Fprintf(outfile, "Finished at time %3d\n\n", schedule.duration)

	for i := 0; i < count; i++ {
		fmt.Fprintf(outfile, "%+v wait %3d turnaround %3d\n", processes[i].pid, wait[i], turnaround[i])
	}
}
