package repository

import (
	"codis/models"
	"strconv"

	"github.com/lib/pq"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type NodeRepository struct {
	postgresDatabaseService *PostgresDatabaseService
}

func NewNodeRepository(injector do.Injector) (*NodeRepository, error) {
	r := NodeRepository{
		postgresDatabaseService: do.MustInvoke[*PostgresDatabaseService](injector),
	}

	return &r, nil
}

func (r NodeRepository) Create(node models.Node) (created models.Node, err error) {
	q := `INSERT INTO public.node (workflow_id, type, next_node_id) VALUES ($1, $2, $3) RETURNING *;`
	err = r.postgresDatabaseService.Get(&created, q, node.WorkflowID, node.Type, node.NextNodeID)
	return
}

func (r NodeRepository) CreateMany(nodes []models.Node) (created []models.Node, err error) {
	q := `INSERT INTO public.node (id, workflow_id, type, next_node_id) VALUES `
	values := []interface{}{}
	for i, node := range nodes {
		q += `($` + strconv.Itoa(i+1) + `, $` + strconv.Itoa(i+2) + `, $` + strconv.Itoa(i+3) + `, $` + strconv.Itoa(i+4) + `)`
		if i < len(nodes)-1 {
			q += `, `
		}
		values = append(values, node.ID, node.WorkflowID, node.Type, node.NextNodeID)
	}
	q += ` RETURNING *;`

	err = r.postgresDatabaseService.Db.Select(&created, q, values...)
	if err != nil {
		return nil, oops.Wrap(err)
	}
	return
}

func (r NodeRepository) GetByID(id string) (node models.Node, err error) {
	q := `SELECT * FROM public.node WHERE id = $1 AND deleted_at IS NULL;`
	err = r.postgresDatabaseService.Get(&node, q, id)
	return
}

func (r NodeRepository) ListByWorkflowID(workflowID string) (nodes []models.Node, err error) {
	q := `SELECT * FROM public.node WHERE workflow_id = $1 AND deleted_at IS NULL ORDER BY created_at;`
	err = r.postgresDatabaseService.Db.Select(&nodes, q, workflowID)
	if err != nil {
		return nil, oops.Wrap(err)
	}
	return
}

func (r NodeRepository) Update(node models.Node) (updated models.Node, err error) {
	q := `UPDATE public.node SET type = $1, next_node_id = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL RETURNING *;`
	err = r.postgresDatabaseService.Get(&updated, q, node.Type, node.NextNodeID, node.ID)
	return
}

func (r NodeRepository) Delete(ids []string) (err error) {
	q := `UPDATE public.node SET deleted_at = NOW() WHERE id = ANY($1);`
	err = r.postgresDatabaseService.Exec(q, pq.StringArray(ids))
	return
}
