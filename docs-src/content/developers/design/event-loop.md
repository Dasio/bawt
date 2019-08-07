---
title: "Event Loop"
description: The event loop is designed to keep Bawt's internal state in line as well as pass messages to listeners. 
weight: 10
---

Core to bawt is listening to an event queue that Slack provides in the form of their Real Time Messaging API. We process this queue through an event loop, which takes it where it needs to go in bawt mainly for internal purposes.

{{<mermaid align="left">}}
graph LR;
    A[slack.RTMEvent] -->|Process Event| B(handleRTMEvent)
    BB -.-> C(slack.LatencyReport)
    BB -.-> D(slack.RTMError)
    BB -.-> E(slack.ConnectedEvent)
    BB -.-> F(slack.DisconnectedEvent)
    BB -.-> G(slack.InvalidAuthEvent)
    BB -.-> H(slack.ConnectingEvent)
    BB -.-> I(slack.HelloEvent)
    BB -.-> J(slack.MessageEvent)
    BB -.-> K(slack.UserChangeEvent)
    BB -.-> L(slack.ChannelRenameEvent)
    BB -.-> M(slack.ChannelJoinedEvent)
    BB -.-> N(slack.ChannelCreatedEvent)
    BB -.-> O(slack.ChannelDeletedEvent)
    BB -.-> P(slack.ChannelArchiveEvent)
    BB -.-> Q(slack.GroupRenameEvent)
    BB -.-> R(slack.GroupJoinedEvent)
    BB -.-> S(slack.GroupCreatedEvent)
    BB -.-> T(slack.GroupCloseEvent)
    BB -.-> U(slack.GroupArchiveEvent)
    BB -.-> V(slack.GroupUnarchivedEvent)
    BB -.-> W(slack.IMCreatedEvent)
    BB -.-> X(slack.IMOpenEvent)
    BB -.-> Y(slack.IMCloseEvent)
    BB -.-> Z(slack.AckErrorEvent)
    BB -.-> AA(slack.ConnectionErrorEvent)
    B -->|First| BB{Event Loop}
    B -->|Second| CC(Send to Listeners Loop)

    click J "/bawt/developers/design/event-loop/#slack-messageevent" "Click for more details"
{{< /mermaid >}}

{{% notice note %}}
bawt's ordered execution of the event loop and the listeners is intentional. Many of bawt's plugins assume that bawt's internal state is up to date and this prevents strange race conditions.
{{% /notice %}}

## Events

### slack.MessageEvent

Notice this is not where messages are routed to listeners. The event loop is internal to bawt and only important for keeping bawt's core state inline.

{{<mermaid align="left">}}
graph LR;
    A[slack.MessageEvent] -->|Incoming Message| B{Event Loop}
    B -->|Determine by Type| C(slack.Message)
    C -.->|Normalize Message| D("&Message{}")
    C -.-> E{ev.Msg.SubType}
    D -.-> F(message_changed)
    D -.-> G(channel_topic)
    D -.-> H(channel_purpose)
    C -.->|Verify UserMap| I(bot.Users)
    C -.->|Verify ChannelMap| J(bot.Channels)
    C -.->|Mutate Message| K(applyMentionsMe)
    C -.->|Mutate Message| L(applyFromMe)
    B -->|Dispatch to Listeners| M((Dispatcher))
    M -.-> N(Listener A)
    N -.->|Message Handler| O(bawt.Listener.filterAndDispatchMessage)
    O -.-> Q(bawt.Listener.filterMessage)
    Q -.-> R(bawt.Listener.MessageHandlerFunc)
    N -.->|Event Handler| P(bawt.Listener.EventHandlerFunc)
{{< /mermaid >}}