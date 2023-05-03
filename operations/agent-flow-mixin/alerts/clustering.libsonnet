local alert = import './utils/alert.jsonnet';
local filename = 'clustering.json';

{
  [filename]:
    alert.new(
      "clustering alerts",
      [
      // Cluster not converging.
      alert.newRule(
        'ClusterNotConverged',
        'stddev(cluster_node_peers) != 0',
        'Cluster is not converging'
      ) + alert.withForTime('2m'),

      // Clustering has entered a split brain state
      alert.newRule(
        'ClusterSplitBrain',
        '(sum without (state) (cluster_node_peers)) != count(cluster_node_info)',
        'Cluster nodes have entered a split brain state'
      ) + alert.withForTime('2m'),
    ]),
}
