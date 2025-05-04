package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"example.com/m/pkg/models"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

type ElasticsearchUtil struct {
	elasticClient *elasticsearch.Client
}

func NewElasticsearchUtil() *ElasticsearchUtil {
	client, err := elasticsearch.NewClient(elasticsearch.Config{})
	if err != nil {
		log.Println("Warning: " + err.Error())
	}

	return &ElasticsearchUtil{
		client,
	}
}

func (e *ElasticsearchUtil) CreateIndex(index string) (*esapi.Response, error) {
	settingData := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   2,
			"number_of_replicas": 2,
		},
	}
	settingDataBytes, err := json.Marshal(settingData)
	if err != nil {
		return nil, err
	}

	res, err := e.elasticClient.Indices.Create(index, func(r *esapi.IndicesCreateRequest) {
		r.Body = bytes.NewReader(settingDataBytes)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *ElasticsearchUtil) IndexDocument(document interface{}, index, id string) (*esapi.Response, error) {
	data, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}

	resp, err := e.elasticClient.Index(index, bytes.NewReader(data), func(r *esapi.IndexRequest) {
		r.DocumentID = id
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *ElasticsearchUtil) BlukIndexDocument(documents []models.Item, index string) error {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: e.elasticClient,
		Index:  index,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < len(documents); i++ {
		wg.Add(1)

		id := documents[i].GetID()
		if id == "" {
			return errors.New("error: can't insert document due to document_id")
		}

		data, err := json.Marshal(documents[i])
		if err != nil {
			return err
		}

		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: id,
				Body:       bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					fmt.Printf("[%d] %s %s\n", res.Status, res.Result, item.DocumentID)
					wg.Done()
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s\n", err)
					} else {
						log.Printf("ERROR: %s: %s\n", res.Error.Type, res.Error.Reason)
					}
					wg.Done()
				},
			},
		)
		if err != nil {
			return err
		}
	}

	wg.Wait()
	if err := indexer.Close(context.Background()); err != nil {
		return err
	}

	return nil
}

func (e *ElasticsearchUtil) UpdateDocumentById(index, id string, updatedDocument interface{}) (*esapi.Response, error) {
	if updatedDocument == nil {
		return nil, errors.New("error: updatedDocument is null")
	}

	updates := map[string]interface{}{
		"doc": updatedDocument,
	}
	data, err := json.Marshal(updates)
	if err != nil {
		return nil, err
	}

	res, err := e.elasticClient.Update(index, id, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *ElasticsearchUtil) DeleteDocumentById(index, id string) (*esapi.Response, error) {
	resp, err := e.elasticClient.Delete(index, id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *ElasticsearchUtil) SearchDocuments(index, value string, size uint) ([]interface{}, error) {
	query := map[string]interface{}{
		"size": size,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query": value,
				},
			},
		},
	}
	queryData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resp, err := e.elasticClient.Search(
		e.elasticClient.Search.WithIndex(index),
		e.elasticClient.Search.WithBody(bytes.NewReader(queryData)),
		e.elasticClient.Search.WithTrackTotalHits(true),
		e.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	var r = map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	var docs []interface{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})["_source"]
		docs = append(docs, doc)
	}

	return docs, nil
}
