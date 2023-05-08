package models

type Configuration struct {
	ID            string `bson:"_id,omitempty"`
	NodeID        string `bson:"node_id"`
	ListenerName  string `bson:"listener_name"`
	RouteConfig   string `bson:"route_config"`
	ClusterConfig string `bson:"cluster_config"`
	FilterConfig  string `bson:"filter_config"`
}
