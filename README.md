# arlo-go
![](gopher-arlo.png)
> Go package for interacting with Netgear's Arlo camera system.

---
### Now in Go!
I love Go. That is why I decided to write this library! I am the creator of the first [arlo](https://github.com/jeffreydwalter/arlo) library written in Python.

My goal is to bring parity to the Python version asap. If you know what you're doing in Go, I would appreciate any feedback on the general structure of the library, bugs found, contributions, etc.

---
It is by no means complete, although it does expose quite a bit of the Arlo interface in an easy to use Go pacakge. As such, this package does not come with unit tests (feel free to add them, or I will eventually) or guarantees.
**All [contributions](https://github.com/jeffreydwalter/arlo/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) are welcome and appreciated!**

**Please, feel free to [contribute](https://github.com/jeffreydwalter/arlo/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) to this repo or buy Jeff a beer!** [![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=R77B7UXMLA6ML&lc=US&item_name=Jeff%20Needs%20Beer&item_number=buyjeffabeer&currency_code=USD&bn=PP%2dDonationsBF%3abtn_donateCC_LG%2egif%3aNonHosted)

---
### Generous Benefactors (Thank you!)
No beers for Jeff yet! 🍺

---
### Awesomely Smart Contributors (Thank you!)
Just me so far...

If You'd like to make a diffrence in the world and get your name on this most prestegious list, have a look at our [help wanted](https://github.com/jeffreydwalter/arlo/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) section!

---
### Filing an Issue
Please read the [Issue Guidelines and Policies](https://github.com/jeffreydwalter/arlo/wiki/Issue-Guidelines-and-Policies) wiki page **BEFORE** you file an issue. Thanks.

---

## Install
```bash
# Install latest stable package
$ go get github.com/jeffreydwalter/arlo-go

# Note: This package uses the `go dep` package for dependency management. If you plan on contributing to this package, you will be required to use [dep](https://github.com/golang/dep). Setting it up is outside the scope of this README, but if you want to contribute and aren't familiar with `dep`, I'm happy to get you.
```

After installing all of the required libraries, you can import and use this library like so:

```golang
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jeffreydwalter/arlo-go"
)

const (
	USERNAME = "user@example.com"
	PASSWORD = "supersecretpassword"
)

func main() {

	// Instantiating the Arlo object automatically calls Login(), which returns an oAuth token that gets cached.
	// Subsequent successful calls to login will update the oAuth token.
	arlo, err := arlo.Login(USERNAME, PASSWORD)
	if err != nil {
		log.Printf("Failed to login: %s\n", err)
		return
	}
	// At this point you're logged into Arlo.

	now := time.Now()
	start := now.Add(-7 * 24 * time.Hour)

	// Get all of the recordings for a date range.
	library, err := arlo.GetLibrary(start, now)
	if err != nil {
		log.Println(err)
		return
	}

	// We need to wait for all of the recordings to download.
	var wg sync.WaitGroup

	for _, recording := range *library {

		// Let the wait group know about the go routine that we're about to run.
		wg.Add(1)

		// The go func() here makes this script download the files concurrently.
		// If you want to download them serially for some reason, just remove the go func() call.
		go func() {
			fileToWrite, err := os.Create(fmt.Sprintf("downloads/%s_%s.mp4", time.Unix(0, recording.UtcCreatedDate*int64(time.Millisecond)).Format(("2006-01-02_15.04.05")), recording.UniqueId))
            defer fileToWrite.Close()

            if err != nil {
                log.Fatal(err)
            }

			// The videos produced by Arlo are pretty small, even in their longest, best quality settings.
			// DownloadFile() efficiently streams the file from the http.Response.Body directly to a file.
			if err := arlo.DownloadFile(recording.PresignedContentUrl, fileToWrite); err != nil {
				log.Println(err)
			} else {
				log.Printf("Downloaded video %s from %s", recording.CreatedDate, recording.PresignedContentUrl)
			}

			// Mark this go routine as done in the wait group.
			wg.Done()
		}()
	}

	// Wait here until all of the go routines are done.
	wg.Wait()


    // The below example demonstrates how you could delete the cloud recordings after downloading them.
    // Simply uncomment the below code to start using it.

    // Delete all of the videos you just downloaded from the Arlo library.
	// Notice that you can pass the "library" object we got back from the GetLibrary() call.
	/* if err := arlo.BatchDeleteRecordings(library); err != nil {
		log.Println(err)
		return
	} */

	// If we made it here without an exception, then the videos were successfully deleted.
	/* log.Println("Batch deletion of videos completed successfully.") */
}
```

** (coming soon) For more code examples check out the [wiki](https://github.com/jeffreydwalter/arlo-go/wiki)**
