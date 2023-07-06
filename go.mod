module github.com/andersnauman/go-ipam

go 1.20

require (
	github.com/andersnauman/go-ipam/api v0.0.0-20230702220342-5eeb20827180
	github.com/andersnauman/go-ipam/config v0.0.0-20230702220342-5eeb20827180
	github.com/andersnauman/go-ipam/ui v0.0.0-20230702220342-5eeb20827180
	github.com/go-sql-driver/mysql v1.7.1
)

require github.com/andersnauman/go-ipam/misc v0.0.0-20230702220342-5eeb20827180 // indirect

replace github.com/andersnauman/go-ipam => ../go-ipam
