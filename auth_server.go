package main

import (
	"context"

	"github.com/boof/umg/auth"
	pb "github.com/boof/umg/proto"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/services"
)

type AuthServer struct{}

func (*AuthServer) HasDomPerm(ctx context.Context, req *pb.DomPermReq) (*pb.PermRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	return &pb.PermRes{Has: err == nil && services.HasDomPerm(user, req.Domain, req.Action)}, nil
}

func (*AuthServer) HasProdPerm(ctx context.Context, req *pb.ProdPermReq) (*pb.PermRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	return &pb.PermRes{Has: err == nil && services.HasProdPerm(user, req.Domain, req.Product, req.Action)}, nil
}

func (*AuthServer) AddProduct(ctx context.Context, req *pb.AddProdReq) (*pb.AddProdRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	if err != nil || !user.IsAdmin() {
		return &pb.AddProdRes{Done: false, Message: "Unauthorized action, only admins can add product"}, nil
	}

	domain, err := (&domains.Domain{Name: req.Domain}).GetByName()
	if err != nil {
		return &pb.AddProdRes{Done: false, Message: "Invalid domain"}, nil
	}

	product := &products.Product{
		DomainID: domain.ID,
		Name:     req.Product,
	}

	if err := product.Save(); err != nil {
		return &pb.AddProdRes{Done: false, Message: err.Error()}, nil
	}

	return &pb.AddProdRes{Done: true}, nil
}

func (*AuthServer) RemoveProduct(ctx context.Context, req *pb.RemProdReq) (*pb.RemProdRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	if err != nil || !user.IsAdmin() {
		return &pb.RemProdRes{Done: false, Message: "Unauthorized action, only admins can remove product"}, nil
	}

	domain, err := (&domains.Domain{Name: req.Domain}).GetByName()
	if err != nil {
		return &pb.RemProdRes{Done: false, Message: "Invalid domain"}, nil
	}

	product, err := (&products.Product{DomainID: domain.ID, Name: req.Product}).GetByName()
	if err != nil {
		return &pb.RemProdRes{Done: false, Message: "Invalid product"}, nil
	}

	if err := product.RemoveByID(); err != nil {
		return &pb.RemProdRes{Done: false, Message: err.Error()}, nil
	}

	return &pb.RemProdRes{Done: true}, nil
}

func (*AuthServer) AddProperty(ctx context.Context, req *pb.AddPropertyReq) (*pb.AddPropertyRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	if err != nil || !user.IsAdmin() {
		return &pb.AddPropertyRes{Done: false, Message: "Unauthorized action, only admins can add property"}, nil
	}

	property := &properties.Property{
		MeteringID: req.Id,
		Type:       req.Type,
		Name:       req.Name,
	}

	err = property.Save()
	if err != nil {
		return &pb.AddPropertyRes{Done: false, Message: err.Error()}, nil
	}

	return &pb.AddPropertyRes{Done: true}, nil
}

func (*AuthServer) RemoveProperty(ctx context.Context, req *pb.RemPropertyReq) (*pb.RemPropertyRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	if err != nil || !user.IsAdmin() {
		return &pb.RemPropertyRes{Done: false, Message: "Unauthorized action, only admins can remove property"}, nil
	}

	property, err := (&properties.Property{MeteringID: req.Id, Type: req.Type}).GetByID()
	if err != nil {
		return &pb.RemPropertyRes{Done: false, Message: "Invalid property"}, nil
	}

	if err := property.RemoveByID(); err != nil {
		return &pb.RemPropertyRes{Done: false, Message: err.Error()}, nil
	}

	return &pb.RemPropertyRes{Done: true}, nil
}

func (*AuthServer) HasPropertyPerm(ctx context.Context, req *pb.PropertyPermReq) (*pb.PermRes, error) {
	user, err := auth.GetUserFromToken(req.Token)
	return &pb.PermRes{Has: err == nil && services.HasPropertyPerm(user, req.Id, req.Type)}, nil
}
