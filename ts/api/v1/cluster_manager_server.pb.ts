/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
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