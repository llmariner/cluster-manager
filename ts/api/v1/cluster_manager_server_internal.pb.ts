/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as LlmarinerClustersServerV1Cluster_manager_server from "./cluster_manager_server.pb"
export type InternalCluster = {
  cluster?: LlmarinerClustersServerV1Cluster_manager_server.Cluster
  tenant_id?: string
}

export type ListInternalClustersRequest = {
}

export type ListInternalClustersResponse = {
  clusters?: InternalCluster[]
}

export class ClustersInternalService {
  static ListInternalClusters(req: ListInternalClustersRequest, initReq?: fm.InitReq): Promise<ListInternalClustersResponse> {
    return fm.fetchReq<ListInternalClustersRequest, ListInternalClustersResponse>(`/llmariner.clusters.server.v1.ClustersInternalService/ListInternalClusters`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}