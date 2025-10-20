package web

import "github.com/curtisnewbie/miso/miso"

// Do not modify this if misoapi is used, misoapi may modify func body
func PrepareWebServer(rail miso.Rail) error {
	RegisterApi()
	RegisterPostboxRoutes(rail)
	return nil
}
