package aws

import (
	"context"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// IAMUser represents an IAM user.
type IAMUser struct {
	UserName         string
	UserID           string
	ARN              string
	CreateDate       string
	PasswordLastUsed string
	Groups           []string
	Policies         []string // attached policy names
}

// CallerIdentity represents the current AWS identity.
type CallerIdentity struct {
	Account string
	ARN     string
	UserID  string
}

// FetchIAMUsers returns all IAM users with their groups and attached policies.
func FetchIAMUsers(ctx context.Context, iamClient *iam.Client) ([]IAMUser, error) {
	paginator := iam.NewListUsersPaginator(iamClient, &iam.ListUsersInput{})
	var users []IAMUser
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, u := range page.Users {
			userName := sdkaws.ToString(u.UserName)

			createDate := ""
			if u.CreateDate != nil {
				createDate = u.CreateDate.Format("2006-01-02 15:04")
			}

			passwordLastUsed := ""
			if u.PasswordLastUsed != nil {
				passwordLastUsed = u.PasswordLastUsed.Format("2006-01-02 15:04")
			}

			// Fetch groups for user
			var groups []string
			groupsOut, err := iamClient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
				UserName: sdkaws.String(userName),
			})
			if err == nil {
				for _, g := range groupsOut.Groups {
					groups = append(groups, sdkaws.ToString(g.GroupName))
				}
			}

			// Fetch attached policies for user
			var policies []string
			policiesOut, err := iamClient.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
				UserName: sdkaws.String(userName),
			})
			if err == nil {
				for _, p := range policiesOut.AttachedPolicies {
					policies = append(policies, sdkaws.ToString(p.PolicyName))
				}
			}

			users = append(users, IAMUser{
				UserName:         userName,
				UserID:           sdkaws.ToString(u.UserId),
				ARN:              sdkaws.ToString(u.Arn),
				CreateDate:       createDate,
				PasswordLastUsed: passwordLastUsed,
				Groups:           groups,
				Policies:         policies,
			})
		}
	}
	return users, nil
}

// FetchCurrentIdentity returns the current caller identity via STS.
func FetchCurrentIdentity(ctx context.Context, stsClient *sts.Client) (*CallerIdentity, error) {
	out, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return &CallerIdentity{
		Account: sdkaws.ToString(out.Account),
		ARN:     sdkaws.ToString(out.Arn),
		UserID:  sdkaws.ToString(out.UserId),
	}, nil
}
