// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIBMSchematicsResourceQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsResourceQueryRead,

		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource query Id.  Use `GET /v2/resource_query` API to look up the Resource query definition Ids  in your IBM Cloud account.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource type (cluster, vsi, icd, vpc).",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource query name.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource Query id.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource query creation time.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who created the Resource query.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource query updation time.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who updated the Resource query.",
			},
			"queries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the query(workspaces).",
						},
						"query_condition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the resource query param.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value of the resource query param.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of resource query param variable.",
									},
								},
							},
						},
						"query_select": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of query selection parameters.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSchematicsResourceQueryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getResourcesQueryOptions := &schematicsv1.GetResourcesQueryOptions{}

	getResourcesQueryOptions.SetQueryID(d.Get("query_id").(string))

	resourceQueryRecord, response, err := schematicsClient.GetResourcesQueryWithContext(context, getResourcesQueryOptions)
	if err != nil {
		log.Printf("[DEBUG] GetResourcesQueryWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetResourcesQueryWithContext failed %s\n%s", err, response))
	}

	d.SetId(*getResourcesQueryOptions.QueryID)
	if err = d.Set("type", resourceQueryRecord.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}
	if err = d.Set("name", resourceQueryRecord.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("id", resourceQueryRecord.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting id: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(resourceQueryRecord.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", resourceQueryRecord.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(resourceQueryRecord.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", resourceQueryRecord.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_by: %s", err))
	}

	if resourceQueryRecord.Queries != nil {
		err = d.Set("queries", dataSourceResourceQueryRecordFlattenQueries(resourceQueryRecord.Queries))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting queries %s", err))
		}
	}

	return nil
}

func dataSourceResourceQueryRecordFlattenQueries(result []schematicsv1.ResourceQuery) (queries []map[string]interface{}) {
	for _, queriesItem := range result {
		queries = append(queries, dataSourceResourceQueryRecordQueriesToMap(queriesItem))
	}

	return queries
}

func dataSourceResourceQueryRecordQueriesToMap(queriesItem schematicsv1.ResourceQuery) (queriesMap map[string]interface{}) {
	queriesMap = map[string]interface{}{}

	if queriesItem.QueryType != nil {
		queriesMap["query_type"] = queriesItem.QueryType
	}
	if queriesItem.QueryCondition != nil {
		queryConditionList := []map[string]interface{}{}
		for _, queryConditionItem := range queriesItem.QueryCondition {
			queryConditionList = append(queryConditionList, dataSourceResourceQueryRecordQueriesQueryConditionToMap(queryConditionItem))
		}
		queriesMap["query_condition"] = queryConditionList
	}
	if queriesItem.QuerySelect != nil {
		queriesMap["query_select"] = queriesItem.QuerySelect
	}

	return queriesMap
}

func dataSourceResourceQueryRecordQueriesQueryConditionToMap(queryConditionItem schematicsv1.ResourceQueryParam) (queryConditionMap map[string]interface{}) {
	queryConditionMap = map[string]interface{}{}

	if queryConditionItem.Name != nil {
		queryConditionMap["name"] = queryConditionItem.Name
	}
	if queryConditionItem.Value != nil {
		queryConditionMap["value"] = queryConditionItem.Value
	}
	if queryConditionItem.Description != nil {
		queryConditionMap["description"] = queryConditionItem.Description
	}

	return queryConditionMap
}
