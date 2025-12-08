package influx

import (
	"context"
	"fmt"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
	"time"

	"log"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type ElectricityRepo struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewElectricityRepo(client influxdb2.Client, org, bucket string) repository.InfluxRepo {
	return &ElectricityRepo{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

func (r *ElectricityRepo) GetRealTimeElectricity(ctx context.Context, deviceID string) (*[]entity.RealTimeElectricity, error) {
	var data entity.RealTimeElectricity

	queryAPI := r.client.QueryAPI(r.org)
	query := fmt.Sprintf(`from(bucket: "%s") 
		|> range(start: -1h) 
		|> filter(fn: (r) => r["_measurement"] == "electricity" and r["device_id"] == "%s") 
		|> sort(columns: ["_time"])
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")`, r.bucket, deviceID)

	res, err := queryAPI.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to query InfluxDB: %w", err)
	}

	var dataList []entity.RealTimeElectricity
	count := 0

	for res.Next() {
		count++
		record := res.Record()

		data.Voltage = record.ValueByKey("voltage").(float64)
		data.Current = record.ValueByKey("current").(float64)
		data.Power = record.ValueByKey("power").(float64)
		data.TotalEnergy = record.ValueByKey("total_energy").(float64)
		data.PowerFactor = record.ValueByKey("power_factor").(float64)
		data.Frequency = record.ValueByKey("frequency").(float64)
		data.DeviceID = record.ValueByKey("device_id").(string)
		data.CreatedAt = utils.NewTimeData(record.Time())

		dataList = append(dataList, data)

	}

	if res.Err() != nil {
		return nil, fmt.Errorf("error reading query result: %w", res.Err())
	}

	if count == 0 {
		log.Printf("Tidak ada data ditemukan untuk device ID: %s", deviceID)
		return nil, nil
	}

	return &dataList, nil
}

func (r *ElectricityRepo) SaveRealTimeElectricity(ctx context.Context, electricity *entity.RealTimeElectricity) error {
	writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint(
		"electricity",
		map[string]string{
			"device_id": electricity.DeviceID,
		},
		map[string]interface{}{
			"voltage":      electricity.Voltage,
			"current":      electricity.Current,
			"power":        electricity.Power,
			"total_energy": electricity.TotalEnergy,
			"power_factor": electricity.PowerFactor,
			"frequency":    electricity.Frequency,
		},
		time.Now().UTC(),
	)

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		return fmt.Errorf("failed to write point to InfluxDB: %w", err)
	}

	return nil
}
