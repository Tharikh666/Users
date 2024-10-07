package user

import (
	"User_Sample/common"
	Users "User_Sample/domain/user"
	"User_Sample/responsecodes"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc"
)

type UserService struct {
	Users.UnimplementedUsersServer
}

func (mObj *UserService) Create(ctx context.Context, Request *Users.User) (*Users.Response, error) {

	defer common.PanicRecovery("UserService", "Create")

	myDB, Ok := common.GetDBConnection()

	if !Ok {

		log.Println("Error in connection")
		return &Users.Response{Id: Request.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, errors.New("error occured in DB connection")
	}
	// SQL query execution
	err := myDB.QueryRow(`INSERT INTO "user" (username,
											password,
											active)
											VALUES ($1,$2,$3) RETURNING id`,
		Request.UserName,
		Request.Password,
		Request.Active).Scan(&Request.Id)

	if err != nil {

		log.Println("UserService", "Create", err)
		return &Users.Response{Id: Request.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, errors.New("query execution failed")
	}

	return &Users.Response{Id: Request.Id,
		Result:         true,
		ResponseStatus: responsecodes.SUCCESS}, nil
}

func (mObj *UserService) GetAll(mEmpty *Users.Empty, stream grpc.ServerStreamingServer[Users.User]) error {

	defer common.PanicRecovery("UserService", "GetAll")

	myDB, Ok := common.GetDBConnection()
	if !Ok {
		log.Println("Error in DB connection")
		return errors.New("failed to connect to DB")
	}

	// SQL query execution
	mRows, err := myDB.Query(`SELECT id, username, password, active FROM "user"`)
	if err != nil {
		log.Println("Query execution failed", err)
		if err == sql.ErrNoRows {
			return errors.New("no rows found")
		}
		return err
	}
	defer mRows.Close()

	// Slice to store query results
	var myResult []*Users.User

	// Iterate through the rows and scan the result into a slice
	for mRows.Next() {
		mItem := &Users.User{}
		err := mRows.Scan(
			&mItem.Id,
			&mItem.UserName,
			&mItem.Password,
			&mItem.Active,
		)
		if err != nil {
			log.Println("Error scanning row", err)
			return err
		}
		myResult = append(myResult, mItem)
	}

	// Check if any error occurred during iteration
	if err = mRows.Err(); err != nil {
		log.Println("Error during row iteration", err)
		return err
	}

	common.PrintJsonFormat("myResult", myResult)
	// Send the result through the stream
	for _, myItem := range myResult {
		if err := stream.Send(myItem); err != nil {
			log.Println("Error sending user data to stream", err)
			return err
		}
	}

	return nil
}

func (mObj *UserService) GetById(ctx context.Context, mId *Users.Identifier) (*Users.User, error) {

	defer common.PanicRecovery("UserService", "GetById")

	myResult := &Users.User{}

	myDB, Ok := common.GetDBConnection()
	if !Ok {
		log.Println("Error in DB connection")
		return myResult, errors.New("failed to connect to DB")
	}

	common.PrintJsonFormat("Id: ", mId.Id)

	// SQL query execution
	err := myDB.QueryRow(`SELECT 	id,
									username,
									password, 
									active 
									FROM "user" 
									WHERE id=$1`, mId.Id).Scan(
		&myResult.Id,
		&myResult.UserName,
		&myResult.Password,
		&myResult.Active)

	if err != nil {
		log.Println("Query execution failed", err)
		if err == sql.ErrNoRows {
			return myResult, errors.New("no rows found")
		}
		return myResult, err
	}

	common.PrintJsonFormat("User: ", myResult)

	return myResult, nil
}

func (mObj *UserService) Update(ctx context.Context, Request *Users.User) (*Users.Response, error) {

	defer common.PanicRecovery("UserService", "Update")

	myDB, Ok := common.GetDBConnection()
	if !Ok {
		log.Println("Error in DB connection")
		return &Users.Response{Id: Request.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, errors.New("failed to connect to DB")
	}

	myUser, err := mObj.AllocatePayload(ctx, Request)
	if err != nil {

		log.Println("Error occured in payload allocation", err)
		return &Users.Response{Id: Request.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, err
	}

	myRows, err := myDB.Query(`UPDATE "user" SET 	username=$1,
													password=$2,
													active=$3
											WHERE	id=$4`,
		myUser.UserName,
		myUser.Password,
		myUser.Active,
		myUser.Id)

	if err != nil {

		log.Println("Error occured in Updation", err)
		return &Users.Response{Id: Request.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, err
	}

	defer myRows.Close()

	return &Users.Response{Id: Request.Id,
		Result:         true,
		ResponseStatus: responsecodes.SUCCESS}, nil
}

func (mObj *UserService) AllocatePayload(ctx context.Context, mReq *Users.User) (*Users.User, error) {

	defer common.PanicRecovery("UserService", "AllocatePayload")

	myUser, err := mObj.GetById(ctx, &Users.Identifier{Id: mReq.Id})
	if err != nil {

		log.Println("UserService", "AllocatePayload", "GetById", err)
		return &Users.User{}, errors.New("no data found")
	}

	if myUser.UserName != mReq.UserName && mReq.UserName != "" {
		myUser.UserName = mReq.UserName
	}

	if myUser.Password != mReq.Password && mReq.Password != "" {
		myUser.Password = mReq.Password
	}

	return myUser, nil
}

func (mObj *UserService) Delete(ctx context.Context, mId *Users.Identifier) (*Users.Response, error) {

	defer common.PanicRecovery("UserService", "Delete")

	myDB, Ok := common.GetDBConnection()
	if !Ok {
		log.Println("Error in DB connection")
		return &Users.Response{Id: mId.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, errors.New("failed to connect to DB")
	}

	mRows, err := myDB.Query(`DELETE FROM "user" WHERE id=$1`, mId.Id)

	if err != nil {

		log.Println("Deleting data failed", err)
		return &Users.Response{Id: mId.Id,
			Result:         false,
			ResponseStatus: responsecodes.FAILED}, err
	}

	defer mRows.Close()

	return &Users.Response{Id: mId.Id,
		Result:         true,
		ResponseStatus: responsecodes.SUCCESS}, nil
}
