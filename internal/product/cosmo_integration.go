package product

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mercadola/api/internal/cosmos"
	"github.com/mercadola/api/internal/infrastruture/config"
)

func GetCosmosProductByEan(ean string) (*Product, error) {
	logger := slog.Default()
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying load config: %s", err.Error()))
	}
	url := fmt.Sprintf("%s/gtins/%s", cfg.CosmoUrl, ean)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao criar a requisição: %s", err.Error()))
		return nil, err
	}

	req.Header.Add("X-Cosmos-Token", cfg.CosmoToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao enviar a requisição: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao ler a resposta: %s", err.Error()))
		return nil, err
	}

	var cosmosProduct cosmos.CosmosProduct
	err = json.Unmarshal(body, &cosmosProduct)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao fazer o unmarshal do JSON: %s", err.Error()))
		return nil, err
	}
	product := FormatToProduct(cosmosProduct)

	return product, nil
}
