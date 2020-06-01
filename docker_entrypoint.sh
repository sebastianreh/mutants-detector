ln -fs "/usr/share/zoneinfo/$TIMEZONE" /etc/localtime && echo "$TIMEZONE" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
go run server.go
