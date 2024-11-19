import * as fm from "../../../fetch.pb";
import * as LlmarinerClustersServerV1Cluster_manager_service from "../cluster_manager_service.pb";
export declare class ClustersInternalService {
    static ListInternalClusters(req: LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersRequest, initReq?: fm.InitReq): Promise<LlmarinerClustersServerV1Cluster_manager_service.ListInternalClustersResponse>;
}
export declare class ClustersWorkerService {
    static GetSelfCluster(req: LlmarinerClustersServerV1Cluster_manager_service.GetSelfClusterRequest, initReq?: fm.InitReq): Promise<LlmarinerClustersServerV1Cluster_manager_service.Cluster>;
}
