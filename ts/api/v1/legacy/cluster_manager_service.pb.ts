/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../../fetch.pb"
import * as LlmarinerClustersServerV1Cluster_manager_service from "../cluster_manager_service.pb"
export class ClustersInternalService {
  static ListInternalClusters(req: LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersRequest, initReq?: fm.InitReq): Promise<LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersResponse> {
    return fm.fetchReq<LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersRequest, LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersResponse>(`/llmoperator.clusters.server.v1.ClustersInternalService/ListInternalClusters`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}
export class ClustersWorkerService {
  static GetSelfCluster(req: LlmarinerClustersServerV1Cluster_manager_service.GetSelfClusterRequest, initReq?: fm.InitReq): Promise<LlmarinerClustersServerV1Cluster_manager_service.Cluster> {
    return fm.fetchReq<LlmarinerClustersServerV1Cluster_manager_service.GetSelfClusterRequest, LlmarinerClustersServerV1Cluster_manager_service.Cluster>(`/llmoperator.clusters.server.v1.ClustersWorkerService/GetSelfCluster`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}