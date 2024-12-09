/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufEmpty from "../../google/protobuf/empty.pb"
export type Cluster = {
  id?: string
  name?: string
  registrationKey?: string
  object?: string
  componentStatuses?: {[key: string]: ComponentStatus}
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

export type InternalCluster = {
  cluster?: Cluster
  tenantId?: string
}

export type ListInternalClustersRequest = {
}

export type ListInternalClustersResponse = {
  clusters?: InternalCluster[]
}

export type GetSelfClusterRequest = {
}

export type ComponentStatus = {
  isHealthy?: boolean
  message?: string
}

export type UpdateComponentStatusRequest = {
  name?: string
  status?: ComponentStatus
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
}
export class ClustersInternalService {
  static ListInternalClusters(req: ListInternalClustersRequest, initReq?: fm.InitReq): Promise<ListInternalClustersResponse> {
    return fm.fetchReq<ListInternalClustersRequest, ListInternalClustersResponse>(`/llmariner.clusters.server.v1.ClustersInternalService/ListInternalClusters`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}
export class ClustersWorkerService {
  static GetSelfCluster(req: GetSelfClusterRequest, initReq?: fm.InitReq): Promise<Cluster> {
    return fm.fetchReq<GetSelfClusterRequest, Cluster>(`/llmariner.clusters.server.v1.ClustersWorkerService/GetSelfCluster`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
  static UpdateComponentStatus(req: UpdateComponentStatusRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<UpdateComponentStatusRequest, GoogleProtobufEmpty.Empty>(`/llmariner.clusters.server.v1.ClustersWorkerService/UpdateComponentStatus`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}