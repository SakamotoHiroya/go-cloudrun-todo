package handler

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	appdb "github.com/SakamotoHiroya/go-cloudrun-todo/db"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/api"
)

func userIDFromContext(ctx context.Context) (int64, error) {
	u, ok := ctx.Value(userKey).(UserInfo)
	if !ok {
		return 0, errors.New("no user in context")
	}
	return strconv.ParseInt(u.id, 10, 64)
}

func apiTaskFromRow(t appdb.Task) api.Task {
	out := api.Task{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.IsCompleted,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Description.Valid {
		out.Description.SetTo(t.Description.String)
	} else {
		out.Description.SetToNull()
	}
	return out
}

func (h *Handler) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (api.CreateTaskRes, error) {
	uid, err := userIDFromContext(ctx)
	if err != nil {
		return &api.Unauthorized{Message: "unauthorized"}, nil
	}

	desc := sql.NullString{}
	if v, ok := req.GetDescription().Get(); ok {
		desc = sql.NullString{String: v, Valid: true}
	}

	row, err := h.repo.CreateTask(ctx, appdb.CreateTaskParams{
		UserID:      uid,
		Title:       req.GetTitle(),
		Description: desc,
	})
	if err != nil {
		return nil, err
	}
	t := apiTaskFromRow(row)
	return &t, nil
}

func (h *Handler) ListTasks(ctx context.Context) (api.ListTasksRes, error) {
	uid, err := userIDFromContext(ctx)
	if err != nil {
		return &api.Unauthorized{Message: "unauthorized"}, nil
	}

	rows, err := h.repo.GetTasks(ctx, uid)
	if err != nil {
		return nil, err
	}
	tasks := make([]api.Task, 0, len(rows))
	for _, r := range rows {
		tasks = append(tasks, apiTaskFromRow(r))
	}
	return &api.TaskListResponse{Tasks: tasks}, nil
}

func (h *Handler) UpdateTaskStatus(ctx context.Context, req *api.UpdateTaskStatusRequest, params api.UpdateTaskStatusParams) (api.UpdateTaskStatusRes, error) {
	uid, err := userIDFromContext(ctx)
	if err != nil {
		return &api.Unauthorized{Message: "unauthorized"}, nil
	}

	row, err := h.repo.UpdateTaskStatus(ctx, appdb.UpdateTaskStatusParams{
		ID:          params.TaskId,
		IsCompleted: req.GetCompleted(),
		UserID:      uid,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nf := api.UpdateTaskStatusNotFound(api.ErrorResponse{Message: "task not found"})
			return &nf, nil
		}
		return nil, err
	}
	t := apiTaskFromRow(row)
	return &t, nil
}

func (h *Handler) DeleteTask(ctx context.Context, params api.DeleteTaskParams) (api.DeleteTaskRes, error) {
	uid, err := userIDFromContext(ctx)
	if err != nil {
		return &api.Unauthorized{Message: "unauthorized"}, nil
	}

	n, err := h.repo.DeleteTask(ctx, appdb.DeleteTaskParams{
		ID:     params.TaskId,
		UserID: uid,
	})
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return &api.ErrorResponse{Message: "task not found"}, nil
	}
	return &api.DeleteTaskNoContent{}, nil
}
