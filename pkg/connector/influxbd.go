package connector

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2writer "github.com/influxdata/influxdb-client-go/v2/api/write"
)

// InfluxDBPoint est un alias pour le type Point de la librairie InfluxDB
type InfluxDBPoint = influxdb2writer.Point

// InfluxService est un service permettant de communiquer avec une base de données InfluxDB.
type InfluxService struct {
	client influxdb2.Client
	org    string
}

// NewInfluxService crée un nouveau service InfluxDB.
func NewInfluxService(url, token, org string) (*InfluxService, error) {
	client := influxdb2.NewClient(url, token)
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, err
	}

	return &InfluxService{
		client: client,
		org:    org,
	}, nil
}

// DataPoint est une interface permettant de convertir un type en un point InfluxDB.
type DataPoint interface {
	ToInfluxDBPoint() *InfluxDBPoint
}

// QueryBuilder est une interface permettant de construire une requête InfluxDB.
type QueryBuilder interface {
	Query() string
}

// Write permet d'écrire un point dans une base de données InfluxDB.
// Le point est écrit dans le bucket spécifié. Si le bucket n'existe pas, il est créé.
func (i *InfluxService) Write(bucket string, data DataPoint) {
	writeAPI := i.client.WriteAPI(i.org, bucket)
	writeAPI.WritePoint(data.ToInfluxDBPoint())
	writeAPI.Flush()
}

// Query permet d'exécuter une requête InfluxDB. Retourne un tableau de map contenant les résultats.
// Exemple de requête : "from(bucket: \"my-bucket\") |> range(start: -1h)"
// Exemple de résultat : [{"_time": "2021-01-01T00:00:00Z", "_value": 42}]
func (i *InfluxService) Query(builder QueryBuilder) ([]map[string]interface{}, error) {
	query := builder.Query()
	qAPI := i.client.QueryAPI(i.org)
	result, err := qAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	defer result.Close()

	output := make([]map[string]interface{}, 0)

	for result.Next() {
		output = append(output, result.Record().Values())
	}

	return output, nil
}
