package product

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/mercadola/api/internal/cosmos"
	"github.com/mercadola/api/internal/infrastruture/config"
)

func GetCosmosProductByEan(ean string) (*Product, error) {
	logger := slog.Default()
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("Error trying load config", err)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/gtins/%s", cfg.CosmoUrl, ean)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Erro ao criar a requisição:", err)
		return nil, err
	}

	req.Header.Add("X-Cosmos-Token", cfg.CosmoToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Erro ao enviar a requisição:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Erro ao ler a resposta:", err)
		return nil, err
	}

	var cosmosProduct cosmos.CosmosProduct
	err = json.Unmarshal(body, &cosmosProduct)
	if err != nil {
		logger.Error("Erro ao fazer o unmarshal do JSON:", err)
		return nil, err
	}
	product := FormatToProduct(cosmosProduct)

	return product, nil
}
