#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { MireiButtonStack } from "../lib/mirei-button-stack";
import { getStage } from "../lib/stage";

const app = new cdk.App();
const stage = getStage(app);
new MireiButtonStack(app, `${stage.stackPrefix}MireiButtonStack`);
