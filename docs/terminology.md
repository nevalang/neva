# Terminology

This document explains key terms used throughout the documentation and how they are applied.

## Paradigm

A programming paradigm is a set of abstractions that describe computation and form a hierarchical tree of different approaches.

```mermaid
---
title: Simplified Overview
---
graph TD
    paradigm --> control-flow
    paradigm --> dataflow
    control-flow --> structural
    control-flow --> object-oriented
    control-flow --> functional
    dataflow --> flow-based
    dataflow --> csp
    dataflow --> actor-model
```

## Control-Flow

Control-flow is a top-level paradigm that describes computation as a series of steps that are executed sequentially.

## Dataflow

Dataflow is a top-level paradigm that describes computation as a network of nodes that perform message-passing.
