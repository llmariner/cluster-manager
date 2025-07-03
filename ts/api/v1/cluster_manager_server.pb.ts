/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufEmpty from "../../google/protobuf/empty.pb"
export type ComponentStatus = {
  is_healthy?: boolean
  message?: string
}

export type Cluster = {
  id?: string
  name?: string
  registration_key?: string
  object?: string
  component_statuses?: {[key: string]: ComponentStatus}
}

export type CreateClusterRequest = {
  name?: string
}

export type ListClustersRequest = {
}

export type ListClustersResponse = {
  object?: string
  data?: Cluster[]
}

export type GetClusterRequest = {
  id?: string
}

export type DeleteClusterRequest = {
  id?: string
}

export type DeleteClusterResponse = {
  id?: string
  object?: string
  deleted?: boolean
}

export type DevicePluginConfigTimeSlicing = {
  gpus?: number
}

export type DevicePluginConfig = {
  time_slicing?: DevicePluginConfigTimeSlicing
}

export type ClusterConfig = {
  device_plugin_config?: DevicePluginConfig
}

export type CreateClusterConfigRequest = {
  cluster_id?: string
  device_plugin_config?: DevicePluginConfig
}

export type GetClusterConfigRequest = {
  cluster_id?: string
}

export type UpdateClusterConfigRequest = {
  cluster_id?: string
  device_plugin_config?: DevicePluginConfig
}

export type DeleteClusterConfigRequest = {
  cluster_id?: string
}

export class ClustersService {
  static CreateCluster(req: CreateClusterRequest, initReq?: fm.InitReq): Promise<Cluster> {
    return fm.fetchReq<CreateClusterRequest, Cluster>(`/v1/clusters`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
  static ListClusters(req: ListClustersRequest, initReq?: fm.InitReq): Promise<ListClustersResponse> {
    return fm.fetchReq<ListClustersRequest, ListClustersResponse>(`/v1/clusters?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GetCluster(req: GetClusterRequest, initReq?: fm.InitReq): Promise<Cluster> {
    return fm.fetchReq<GetClusterRequest, Cluster>(`/v1/clusters/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static DeleteCluster(req: DeleteClusterRequest, initReq?: fm.InitReq): Promise<DeleteClusterResponse> {
    return fm.fetchReq<DeleteClusterRequest, DeleteClusterResponse>(`/v1/clusters/${req["id"]}`, {...initReq, method: "DELETE"})
  }
  static CreateClusterConfig(req: CreateClusterConfigRequest, initReq?: fm.InitReq): Promise<ClusterConfig> {
    return fm.fetchReq<CreateClusterConfigRequest, ClusterConfig>(`/v1/clusters/${req["cluster_id"]}/config`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
  static GetClusterConfig(req: GetClusterConfigRequest, initReq?: fm.InitReq): Promise<ClusterConfig> {
    return fm.fetchReq<GetClusterConfigRequest, ClusterConfig>(`/v1/clusters/${req["cluster_id"]}/config?${fm.renderURLSearchParams(req, ["cluster_id"])}`, {...initReq, method: "GET"})
  }
  static UpdateClusterConfig(req: UpdateClusterConfigRequest, initReq?: fm.InitReq): Promise<ClusterConfig> {
    return fm.fetchReq<UpdateClusterConfigRequest, ClusterConfig>(`/v1/clusters/${req["cluster_id"]}/config`, {...initReq, method: "PATCH", body: JSON.stringify(req)})
  }
  static DeleteClusterConfig(req: DeleteClusterConfigRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<DeleteClusterConfigRequest, GoogleProtobufEmpty.Empty>(`/v1/clusters/${req["cluster_id"]}/config`, {...initReq, method: "DELETE"})
  }
}