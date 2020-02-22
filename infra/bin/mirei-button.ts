#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { MireiButtonStack } from "../lib/mirei-button-stack";

const app = new cdk.App();
new MireiButtonStack(app, "MireiButtonStack");
