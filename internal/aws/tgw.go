package aws

import (
	"context"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TransitGateway represents a TGW.
type TransitGateway struct {
	ID          string
	Name        string
	State       string
	OwnerID     string
	ASN         int64
	CIDR        []string
	Attachments []TGWAttachmentDetail
	RouteTables []TGWRouteTable
}

// TGWAttachmentDetail represents a TGW attachment.
type TGWAttachmentDetail struct {
	ID           string
	ResourceType string
	ResourceID   string
	State        string
}

// TGWRouteTable represents a TGW route table.
type TGWRouteTable struct {
	ID     string
	Name   string
	Routes []TGWRoute
}

// TGWRoute represents a route in a TGW route table.
type TGWRoute struct {
	DestCIDR     string
	AttachmentID string
	ResourceType string
	State        string
}

// FetchTransitGateways returns all transit gateways.
func FetchTransitGateways(ctx context.Context, ec2Client *ec2.Client) ([]TransitGateway, error) {
	paginator := ec2.NewDescribeTransitGatewaysPaginator(ec2Client, &ec2.DescribeTransitGatewaysInput{})
	var gateways []TransitGateway
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, tgw := range page.TransitGateways {
			name := ""
			for _, tag := range tgw.Tags {
				if sdkaws.ToString(tag.Key) == "Name" {
					name = sdkaws.ToString(tag.Value)
					break
				}
			}

			var asn int64
			if tgw.Options != nil && tgw.Options.AmazonSideAsn != nil {
				asn = sdkaws.ToInt64(tgw.Options.AmazonSideAsn)
			}

			var cidrs []string
			if tgw.Options != nil {
				cidrs = tgw.Options.TransitGatewayCidrBlocks
			}

			gateways = append(gateways, TransitGateway{
				ID:      sdkaws.ToString(tgw.TransitGatewayId),
				Name:    name,
				State:   string(tgw.State),
				OwnerID: sdkaws.ToString(tgw.OwnerId),
				ASN:     asn,
				CIDR:    cidrs,
			})
		}
	}
	return gateways, nil
}

// FetchTGWAttachmentsForGateway returns attachments for a specific TGW.
func FetchTGWAttachmentsForGateway(ctx context.Context, ec2Client *ec2.Client, tgwID string) ([]TGWAttachmentDetail, error) {
	paginator := ec2.NewDescribeTransitGatewayAttachmentsPaginator(ec2Client, &ec2.DescribeTransitGatewayAttachmentsInput{
		Filters: []ec2types.Filter{
			{
				Name:   sdkaws.String("transit-gateway-id"),
				Values: []string{tgwID},
			},
		},
	})
	var attachments []TGWAttachmentDetail
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, att := range page.TransitGatewayAttachments {
			attachments = append(attachments, TGWAttachmentDetail{
				ID:           sdkaws.ToString(att.TransitGatewayAttachmentId),
				ResourceType: string(att.ResourceType),
				ResourceID:   sdkaws.ToString(att.ResourceId),
				State:        string(att.State),
			})
		}
	}
	return attachments, nil
}

// FetchTGWRouteTables returns route tables and their routes for a specific TGW.
func FetchTGWRouteTables(ctx context.Context, ec2Client *ec2.Client, tgwID string) ([]TGWRouteTable, error) {
	rtOut, err := ec2Client.DescribeTransitGatewayRouteTables(ctx, &ec2.DescribeTransitGatewayRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   sdkaws.String("transit-gateway-id"),
				Values: []string{tgwID},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var tables []TGWRouteTable
	for _, rt := range rtOut.TransitGatewayRouteTables {
		name := ""
		for _, tag := range rt.Tags {
			if sdkaws.ToString(tag.Key) == "Name" {
				name = sdkaws.ToString(tag.Value)
				break
			}
		}

		rtID := sdkaws.ToString(rt.TransitGatewayRouteTableId)

		// Fetch routes for this route table
		routesOut, err := ec2Client.SearchTransitGatewayRoutes(ctx, &ec2.SearchTransitGatewayRoutesInput{
			TransitGatewayRouteTableId: sdkaws.String(rtID),
			Filters: []ec2types.Filter{
				{
					Name:   sdkaws.String("type"),
					Values: []string{"static", "propagated"},
				},
			},
		})

		var routes []TGWRoute
		if err == nil {
			for _, route := range routesOut.Routes {
				attachmentID := ""
				resourceType := ""
				if len(route.TransitGatewayAttachments) > 0 {
					attachmentID = sdkaws.ToString(route.TransitGatewayAttachments[0].TransitGatewayAttachmentId)
					resourceType = string(route.TransitGatewayAttachments[0].ResourceType)
				}
				routes = append(routes, TGWRoute{
					DestCIDR:     sdkaws.ToString(route.DestinationCidrBlock),
					AttachmentID: attachmentID,
					ResourceType: resourceType,
					State:        string(route.State),
				})
			}
		}

		tables = append(tables, TGWRouteTable{
			ID:     rtID,
			Name:   name,
			Routes: routes,
		})
	}
	return tables, nil
}
