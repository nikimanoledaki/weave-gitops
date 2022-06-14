import * as React from "react";
import styled from "styled-components";
import Interval from "../components/Interval";
import Link from "../components/Link";
import SourceDetail from "../components/SourceDetail";
import Timestamp from "../components/Timestamp";
import { removeKind } from "../lib/utils";
import { FluxObjectKind, HelmRepository } from "../lib/api/core/types.pb";

type Props = {
  className?: string;
  name: string;
  namespace: string;
  clusterName: string;
};

function HelmRepositoryDetail({
  name,
  namespace,
  className,
  clusterName,
}: Props) {
  return (
    <SourceDetail
      className={className}
      name={name}
      namespace={namespace}
      clusterName={clusterName}
      type={FluxObjectKind.KindHelmRepository}
      // Guard against an undefined repo with a default empty object
      info={(hr: HelmRepository = {}) => [
        ["Type", removeKind(FluxObjectKind.KindHelmRepository)],
        ["Repository Type", hr.repositoryType.toLowerCase()],
        ["URL", tryLink(hr.url)],
        [
          "Last Updated",
          hr.lastUpdatedAt ? <Timestamp time={hr.lastUpdatedAt} /> : "-",
        ],
        ["Interval", <Interval interval={hr.interval} />],
        ["Cluster", hr.clusterName],
        ["Namespace", hr.namespace],
      ]}
    />
  );
}

function tryLink(url: string) {
  if (url.startsWith("oci://")) {
    return url;
  } else {
    return (
      <Link newTab href={url}>
        {url}
      </Link>
    );
  }
}

export default styled(HelmRepositoryDetail).attrs({
  className: HelmRepositoryDetail.name,
})``;
