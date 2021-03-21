deploy:
	gcloud functions deploy --project mcg-invitation InviteAllMCG --runtime go113 --trigger-http
	gcloud functions deploy --project mcg-invitation InvitePubSub --trigger-topic triger-invbot --runtime go113