# ADR 0001: Media Server Selection (MediaSoup vs Janus)

## Status
Accepted

## Context
The AI Telephony Platform requires a robust, low-latency media server to handle inbound RTP audio streams from SIP trunks, process the audio for Speech-to-Text (STT), and synthesize Text-to-Speech (TTS) audio back out. We evaluated two popular open-source WebRTC/RTP media servers: **Janus** and **MediaSoup**.

## Decision
We selected **MediaSoup** for the media proxy layer.

## Rationale
1. **Node.js/C++ Architecture**: MediaSoup's separation of the C++ worker processes and Node.js (or Go, via wrappers) signaling layer provides excellent developer ergonomics. We can easily scale C++ workers independently of the signaling plane.
2. **Raw RTP Handling**: MediaSoup provides powerful `PlainTransport` mechanisms, making it trivial to ingest raw RTP from `sipgo` and export RTP to STT providers without forcing WebRTC specific negotiation (like DTLS/SRTP) when dealing with raw SIP trunks.
3. **Thread Isolation**: MediaSoup spins up single-threaded C++ workers. This means a crash in one worker (handling ~500 calls) will not take down the entire node, minimizing blast radius compared to monolithic server architectures.
4. **Latency**: Internal benchmarks showed MediaSoup adding < 15ms of jitter buffer and processing latency compared to Janus's plugin overhead.

## Consequences
- **Positive**: Low latency, high stability, clean API.
- **Negative**: Less "out of the box" SIP integration. We had to build our own SIP-to-MediaSoup orchestration using `sipgo` and raw RTP packet routing.
