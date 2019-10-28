// cleanup allows you to define a cleanup function that will be executed
// if your test is interrupted.

package test

import (
	"os"
	"os/signal"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func CleanupOnInterrupt(cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanup()
			os.Exit(1)
		}
	}()
}

// TearDown will delete created names using clients.
func TearDown(clients *Clients, names ResourceNames) {
	if clients != nil && clients.Eventing != nil {
		clients.KnativeEventing().Delete(names.KnativeEventing, &metav1.DeleteOptions{})
	}
}
