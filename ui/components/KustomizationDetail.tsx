import * as React from "react";
import styled from "styled-components";
import { FluxObjectKind, Kustomization } from "../lib/api/core/types.pb";
import { automationLastUpdated } from "../lib/utils";
import Alert from "./Alert";
import AutomationDetail from "./AutomationDetail";
import Interval from "./Interval";
import SourceLink from "./SourceLink";
import Timestamp from "./Timestamp";

type Props = {
  kustomization?: Kustomization;
  className?: string;
};

function KustomizationDetail({ kustomization, className }: Props) {
  return (
    <AutomationDetail
      className={className}
      automation={{
        ...kustomization,
        kind: FluxObjectKind.KindKustomization,
      }}
      info={[
        [
          "Source",
          <SourceLink
            sourceRef={kustomization?.sourceRef}
            clusterName={kustomization?.clusterName}
          />,
        ],
        ["Applied Revision", kustomization?.lastAppliedRevision],
        ["Cluster", kustomization?.clusterName],
        ["Path", kustomization?.path],

        ["Interval", <Interval interval={kustomization?.interval} />],
        [
          "Last Updated",
          <Timestamp time={automationLastUpdated(kustomization)} />,
        ],
      ]}
    />
  );
}

export default styled(KustomizationDetail).attrs({
  className: KustomizationDetail.name,
})`
  width: 100%;

  ${Alert} {
    margin-bottom: 16px;
  }
`;
