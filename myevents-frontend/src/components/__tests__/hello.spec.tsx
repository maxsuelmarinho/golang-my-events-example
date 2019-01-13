import * as React from "react";
import { shallow } from "enzyme";

import {Hello} from "../hello";

it("renders the div", () => {
    const result = shallow(<Hello name="world" />).contains(<div>Hello world!</div>);
    expect(result).toBeTruthy();
});