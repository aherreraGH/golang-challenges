package main

import "localpractice5b.com/challenges/types"

type FakeClient struct{}

func (f *FakeClient) Call(
	serviceMethod string,
	args any,
	reply any,
) error {

	r := reply.(*types.DeploymentResponse)

	r.Success = true
	r.RequestID = "REQ-123"
	r.Message = "Deployment accepted"

	return nil
}

func (f *FakeClient) Close() error {
	return nil
}
