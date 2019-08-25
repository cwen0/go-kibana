package kibana

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDiscover(t *testing.T) {
	cli := NewDiscover(&Config{Address: "http://172.16.4.4:30750"}, http.DefaultClient)

	search := &SearchSource{
		Search: SearchFields{
			Columns: []string{"log"},
			Sort:    []string{"@timestamp", "desc"},
			Query: QueryMeta{
				Language: "lucene",
			},
			Index:    "d575b240-bdc8-11e9-a168-b9df9becdde7",
			Interval: "auto",
			Filters: []Filter{
				{
					State: State{
						Store: "appState",
					},
					Meta: FilterMeta{
						Disable: false,
						Index:   "d575b240-bdc8-11e9-a168-b9df9becdde7",
						Key:     "kubernetes.namespace_name",
						Negate:  false,
						Params: FilterParams{
							Query: "general-master-merge-all-exp149-cat0-tidb-cluster",
							Type:  "phrase",
						},
						Type:  "phrase",
						Value: "general-master-merge-all-exp149-cat0-tidb-cluster",
					},
					Query: FilterQueryMeta{
						Match: map[string]FilterParams{
							"kubernetes.namespace_name": {
								Query: "general-master-merge-all-exp149-cat0-tidb-cluster",
								Type:  "phrase",
							},
						},
					},
				},
				{
					State: State{
						Store: "appState",
					},
					Meta: FilterMeta{
						Disable: false,
						Index:   "d575b240-bdc8-11e9-a168-b9df9becdde7",
						Key:     "kubernetes.pod_name",
						Negate:  false,
						Params: FilterParams{
							Query: "tidb-cluster-tidb-0",
							Type:  "phrase",
						},
						Type:  "phrase",
						Value: "tidb-cluster-tidb-0",
					},
					Query: FilterQueryMeta{
						Match: map[string]FilterParams{
							"kubernetes.pod_name": {
								Query: "tidb-cluster-tidb-0",
								Type:  "phrase",
							},
						},
					},
				},
			},
		},
	}

	var (
		generalURL string
		err        error
	)
	t.Run("GenerateURL", func(t *testing.T) {
		generalURL, err = cli.GenerateURL(search)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(generalURL)
	})

	t.Run("ShortURL", func(t *testing.T) {
		shortURL, err := cli.ShortURL(generalURL)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(shortURL)
	})
}
