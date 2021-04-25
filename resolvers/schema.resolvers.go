package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"

	"github.com/kotojo/life-manager/generated"
	"github.com/kotojo/life-manager/middleware"
	"github.com/kotojo/life-manager/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *documentResolver) User(ctx context.Context, doc *models.Document) (*models.User, error) {
	user, err := r.UsersService.GetBy("id", doc.UserId)
	if err != nil {
		return nil, gqlerror.Errorf("Error getting user")
	}
	return user, nil
}

func (r *mutationResolver) CreateDocument(ctx context.Context, input models.NewDocument) (*models.Document, error) {
	document, err := r.DocumentsService.Create(input)
	if err != nil {
		return nil, gqlerror.Errorf("Error creating document")
	}
	return document, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUser) (*models.User, error) {
	user, err := r.UsersService.Create(input)
	if err != nil {
		return nil, gqlerror.Errorf("Error creating user")
	}
	err = r.signIn(ctx, user)
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in")
	}
	return user, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input models.UserSignIn) (*models.User, error) {
	user, err := r.UsersService.Authenticate(input.Email, input.Password)
	if err == models.ErrInvalidCredentials {
		return nil, gqlerror.Errorf(err.Error())
	}
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in. Please try again")
	}
	err = r.signIn(ctx, user)
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in")
	}
	return user, nil
}

func (r *mutationResolver) signIn(ctx context.Context, user *models.User) error {
	if user.Remember == "" {
		token, err := r.UsersService.GenerateRememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = r.UsersService.Update(user)
		if err != nil {
			return err
		}
	}
	setCookie(ctx, "remember_token", user.Remember)
	return nil
}

func setCookie(ctx context.Context, name, value string) {
	httpContextValues := ctx.Value(middleware.HttpContextKey).(middleware.HttpContextValues)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	}
	http.SetCookie(*httpContextValues.W, &cookie)
}

func (r *queryResolver) Documents(ctx context.Context) ([]*models.Document, error) {
	user := middleware.UserFromContext(ctx)
	if user == nil {
		return nil, gqlerror.Errorf("User not authorized")
	}
	documents, err := r.DocumentsService.GetAll()
	if err != nil {
		return nil, gqlerror.Errorf("Error getting documents")
	}
	return documents, nil
}

// Document returns generated.DocumentResolver implementation.
func (r *Resolver) Document() generated.DocumentResolver { return &documentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type documentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
