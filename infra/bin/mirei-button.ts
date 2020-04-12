#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { MireiButtonStack } from "../lib/mirei-button-stack";
import { MireiButtonTrimmerStack } from "../lib/mirei-button-trimmer-stack";
import { getStage } from "../lib/stage";

const app = new cdk.App();
const stage = getStage(app);
const props = {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION
  }
};
new MireiButtonStack(app, `${stage.stackPrefix}MireiButtonStack`, props);
new MireiButtonTrimmerStack(
  app,
  `${stage.stackPrefix}MireiButtonTrimmerStack`,
  props
);
app.synth();
