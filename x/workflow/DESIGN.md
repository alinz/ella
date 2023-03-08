# Workflow

```rpc
workflow PaymentWorkflow {
  Checkout()
}

activity StripeActivity {
  Session() => (sessionId: String)
}
```

- calling a workflow
- generate an unique id, workflow_1234
- any call to activity, will make an api to engine /rpc/call/activity
  { workflow_id: "workflow_1234", activity_id: "activity_1234", name: "StripeActivity.Session", args: [], status: "called" }
- call activity
- if no error call /rpc/call/activity
  { workflow_id: "workflow_1234", activity_id: "activity_1234", name: "StripeActivity.Session", args: [], status: "succeed", returns: [] }
- if an error happens
  { workflow_id: "workflow_1234", activity_id: "activity_1234", name: "StripeActivity.Session", args: [], status: "failed", error: "there is an error" }

- topics.workflows.id

```golang

rpc.NewWorker()
rpc.NewEngine

engine := workflow.NewEngine()

rpc.RegisterActivity(engine, rpc.NewStripeActivity)
rpc.RegisterWorkflow(engine, rpc.NewPaymentWorkflow)

rpc.Invoke(payment)

func (ctx workflow.Context) {
  rpc.InvokeActivity(stripeActivity, )
}

rpc.InvokeActivity(engine, stripeActivity.)

```
