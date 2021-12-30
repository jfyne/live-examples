module github.com/jfyne/live-examples

go 1.16

require (
	github.com/gofiber/adaptor/v2 v2.1.15
	github.com/gofiber/fiber/v2 v2.23.0
	github.com/jfyne/live v0.12.3
	github.com/jfyne/live-contrib/livefiber v0.0.0-00010101000000-000000000000
	github.com/maragudk/gomponents v0.16.0
	gocloud.dev v0.22.0
)

replace github.com/jfyne/live => ../live

replace github.com/jfyne/live-contrib/livefiber => ../live-contrib/livefiber
