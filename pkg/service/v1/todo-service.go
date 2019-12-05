package v1

import (
	"context"
	"database/sql"
	"fmt"
	v1 "github.com/donggangcj/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoService is implementation of v1.ToDoServiceServer proto interface
type toDoServiceService struct {
	db *sql.DB
}

//NewToDoServiceServer create a ToDo service
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceService{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceService) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unSupport API version: Service implements API version '%s',but asked for '%s'", apiVersion, api)
		}
	}
	return nil

}

// connect return SQL database connection form the pool
func (s *toDoServiceService) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database -> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s toDoServiceService) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// insert ToDo Entity data
	res, err := c.ExecContext(ctx, "INSERT INTO ToDo(`Title`,`Description`,`Reminder`) VALUES (?,?,?)", req.ToDo.Title, req.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert int ToDo->"+err.Error())
	}

	// get Id of creates ToDo
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo->"+err.Error())
	}
	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil

}

// Read todo task
func (s toDoServiceService) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query ToDO by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`,`Title`,`Description`,`Reminder` FROM  ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo->"+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, ":failed to retrieve data from ToDo->"+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID ='%d' is not found ", req.Id))
	}

	// get ToDo data
	var td v1.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row->"+err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with Id='%d'", req.Id))
	}
	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &td,
	}, nil
}

// update todo task
func (s toDoServiceService) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalide format->"+err.Error())
	}
	// updata ToDo
	res, err := c.ExecContext(ctx, "UPDATE ToDo SET `Title`=?,`Description`=?,`Reminder`=? WHERE `ID`=?", req.ToDo.Title, req.ToDo.Description, reminder, req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affect value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.ToDo.Id))
	}
	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

//Delete todo task
func (s toDoServiceService) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo->"+err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}
	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

//ReadAll read all todo task
func (s toDoServiceService) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get ToDo list
	rows, err := c.QueryContext(ctx, "SELECT `ID`,`Title`,`Description`,`Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	var list []*v1.ToDo
	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}
		td.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format->"+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}
	return &v1.ReadAllResponse{Api: apiVersion, ToDos: list}, nil
}
