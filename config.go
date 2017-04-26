package main

import (
	"bytes"
	"os"

	"encoding/json"
)

type configDB struct {
	Name     string
	User     string
	Password string
	Host     string
	SslMode  string
}

type configServer struct {
	Host         string
	HostedDomain string
	Ui           string
	Home         string
	DisableCors  bool
}

type configOauth struct {
	Secret string
	Id     string
}

type configSsl struct {
	Certificate string
	Key         string
}

type configuration struct {
	Db     configDB
	Server configServer
	Oauth  configOauth
	Ssl    configSsl
}

var (
	Configuration = configuration{
		Db: configDB{
			Name: "securidash",
			User: "securidash",
			Host: "localhost",
		},
		Oauth: configOauth{
			Secret: "oauth.client.secret",
			Id:     "oauth.client.id",
		},
		Ssl: configSsl{
			Certificate: "server.crt",
			Key:         "server.key",
		},
		Server: configServer{
			Host:         "localhost",
			HostedDomain: "",
			Ui:           "/ui/",
			Home:         os.Getenv("HOME"),
			DisableCors:  false,
		},
	}
)

func LoadConfig(config string) error {
	file, err := os.Open(config)
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&Configuration)
	if err != nil {
		return err
	}

	if Configuration.Server.Ui[0] != '/' {
		Configuration.Server.Ui = "/" + Configuration.Server.Ui
	}

	if Configuration.Server.Ui[len(Configuration.Server.Ui)-1] != '/' {
		Configuration.Server.Ui += "/"
	}

	if Configuration.Server.Home[len(Configuration.Server.Home)-1] == '/' {
		Configuration.Server.Home = Configuration.Server.Home[:len(Configuration.Server.Home)-1]
	}

	return nil
}

func GenerateDefaultConfig(file string) error {
	newFile, err := os.Create(file)
	if err != nil {
		return err
	}

	config, err := json.Marshal(Configuration)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = json.Indent(&buffer, config, "", "\t")
	if err != nil {
		return err
	}

	_, err = buffer.WriteTo(newFile)
	if err != nil {
		return err
	}

	return nil
}
