import {
  expect as expectCDK,
  matchTemplate,
  MatchStyle
} from "@aws-cdk/assert";
import * as cdk from "@aws-cdk/core";
import MireiButton = require("../lib/mirei-button-stack");

test("Empty Stack", () => {
  const app = new cdk.App();
  // WHEN
  const stack = new MireiButton.MireiButtonStack(app, "MyTestStack");
  // THEN
  expectCDK(stack).to(
    matchTemplate(
      {
        Resources: {}
      },
      MatchStyle.EXACT
    )
  );
});
