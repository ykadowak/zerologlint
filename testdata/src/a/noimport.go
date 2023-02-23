package a

import "log"

func noimport() {
	log.Fatal("foo") // no report for this because zerolog is not imported
}
