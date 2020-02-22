import * as aws from "aws-sdk";
import * as fs from "fs";
import * as path from "path";
import * as mimeTypes from "mime-types";

const findStack = (): Promise<aws.CloudFormation.Stack> => {
  return new aws.CloudFormation()
    .describeStacks({
      StackName: process.env.MB_STACK_NAME || ""
    })
    .promise()
    .then(stacks => (stacks.Stacks || [])[0]);
};

const findBucketName = (stack: aws.CloudFormation.Stack): string => {
  return (
    stack.Outputs?.find(o => o.OutputKey === "BucketName")?.OutputValue || ""
  );
};

const putObject = (
  bucketName: string,
  basePath: string,
  key: string
): Promise<aws.S3.PutObjectOutput> => {
  const filePath = path.join(basePath, key);
  const file = fs.readFileSync(filePath);

  const contentType = mimeTypes.lookup(filePath) || undefined;

  return new aws.S3()
    .putObject({
      Bucket: bucketName,
      Key: key,
      Body: file,
      ACL: "public-read",
      ContentType: contentType
    })
    .promise();
};

const main = async (): Promise<void> => {
  aws.config.logger = console;

  const stack = await findStack();
  const bucketName = findBucketName(stack);

  await putObject(bucketName, "../public", "index.html");
};

main()
  .then(() => {
    process.exit(0);
  })
  .catch(e => {
    console.error(e);
    process.exit(1);
  });
