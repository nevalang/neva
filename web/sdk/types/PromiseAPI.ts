import { ResponseContext, RequestContext, HttpFile } from '../http/http';
import * as models from '../models/all';
import { Configuration} from '../configuration'

import { Connection } from '../models/Connection';
import { IO } from '../models/IO';
import { Module } from '../models/Module';
import { SDKOperator } from '../models/Operator';
import { PortAddr } from '../models/PortAddr';
import { Program } from '../models/Program';
import { ObservableDefaultApi } from './ObservableAPI';

import { DefaultApiRequestFactory, DefaultApiResponseProcessor} from "../apis/DefaultApi";
export class PromiseDefaultApi {
    private api: ObservableDefaultApi

    public constructor(
        configuration: Configuration,
        requestFactory?: DefaultApiRequestFactory,
        responseProcessor?: DefaultApiResponseProcessor
    ) {
        this.api = new ObservableDefaultApi(configuration, requestFactory, responseProcessor);
    }

    /**
     */
    public programGet(_options?: Configuration): Promise<Program> {
        const result = this.api.programGet(_options);
        return result.toPromise();
    }


}



