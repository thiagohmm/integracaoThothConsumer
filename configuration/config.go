package configuration

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver    string `mapstructure:"DB_DIALECT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWD"`
	DBConnect   string `mapstructure:"DB_CONNECTSTRING"`
	ServiceName string
	Port        int
	Host        string
}

type Dados struct {
	ServiceName string
	Port        int
	Host        string
}

func extrairDados(descricao string) (Dados, error) {
	var dados Dados

	// Expressões regulares para extrair os valores
	serviceNamePattern := regexp.MustCompile(`service_name=([\w.]+)`)
	portPattern := regexp.MustCompile(`port=(\d+)`)
	hostPattern := regexp.MustCompile(`host=([\d.]+)`)

	// Encontrar os valores usando as expressões regulares
	serviceNameMatch := serviceNamePattern.FindStringSubmatch(descricao)
	portMatch := portPattern.FindStringSubmatch(descricao)
	hostMatch := hostPattern.FindStringSubmatch(descricao)

	// Verificar se todos os valores foram encontrados
	if len(serviceNameMatch) < 2 || len(portMatch) < 2 || len(hostMatch) < 2 {
		return dados, fmt.Errorf("não foi possível extrair todos os dados")
	}

	dados.ServiceName = serviceNameMatch[1]
	dados.Port, _ = strconv.Atoi(portMatch[1])
	dados.Host = hostMatch[1]

	return dados, nil
}

// Carrega as configurações do arquivo .env e das variáveis de ambiente
func LoadConfig(path string) (*Conf, error) {
	var cfg Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Não foi possível ler o arquivo .env, tentando variáveis de ambiente")
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	// Extrair os dados da string de conexão
	dados, err := extrairDados(cfg.DBConnect)
	if err != nil {
		return &cfg, err
	}

	// Preencher os campos na struct de configuração
	cfg.ServiceName = dados.ServiceName
	cfg.Port = dados.Port
	cfg.Host = dados.Host

	return &cfg, nil
}
