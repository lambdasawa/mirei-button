import { Construct } from "@aws-cdk/core";

export interface Stage {
  stackPrefix: string;
}

const knownStageNames = ["dev", "prod"];

export const getStage = (construct: Construct): Stage => {
  const stageName = process.env.MB_STAGE || "dev";
  if (!knownStageNames.includes(stageName)) {
    throw new Error(`unknown stage: ${stageName}`);
  }

  const stage = construct.node.tryGetContext(stageName) as Stage;
  if (!stage.stackPrefix) {
    throw new Error(`stack prefix is empty: ${stage.stackPrefix}`);
  }

  console.log({ stage });

  return stage;
};
